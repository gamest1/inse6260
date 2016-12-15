package scheduler

import (
  "fmt"
  "time"
  "math"

  log "github.com/goinggo/tracelog"

  //Database access:
  "github.com/goinggo/beego-mgo/services"
	"github.com/goinggo/beego-mgo/services/userService"
	"github.com/goinggo/beego-mgo/services/requestService"
	"github.com/goinggo/beego-mgo/utilities/mongo"
  "gopkg.in/mgo.v2/bson"
)

type (
  Appointment struct {
    StartTime   time.Time
	  Duration    int
  }

  IMBookings struct {
    Allocations []Appointment
  }

  IMRequest struct {
    ID         bson.ObjectId
    CareGivers []string //Care Givers identified by email!
    Status      string
    CareGiver   string
    TimeInfo    Appointment
  }

  SchedulingSolution struct {
    Solution    []IMRequest
  }

  SolutionBoard struct {
    AllCareGivers map[string]IMBookings
    CurrentSolution *SchedulingSolution
    BestSolution *SchedulingSolution
    IsKilled bool
  }

)

func (a Appointment) Equals(b Appointment) bool {
  if a.StartTime == b.StartTime && a.Duration == b.Duration {
    return true
  }
  return false
}

func (a Appointment) ConflictsWith(b Appointment) bool {
  if a.StartTime.Day()   == b.StartTime.Day()   &&
  	 a.StartTime.Month() == b.StartTime.Month() &&
  	 a.StartTime.Year()  == b.StartTime.Year()  {

    //They are on the same day! There may be a conflict:
    x1 := float64(a.StartTime.Hour())
    x2 := x1 + float64(a.Duration)

    y1 := float64(b.StartTime.Hour())
    y2 := y1 + float64(b.Duration)
    //We just need to find if the intervals [x1,x2] and [y1,y2] intersect!
    return math.Max(x1,y1) < math.Min(x2,y2)
  }

  return false
}

func (b *SolutionBoard) CanPlay(careGiver string, req IMRequest) bool {
  if b.IsKilled {
    return false
  }
  for _, app := range b.AllCareGivers[careGiver].Allocations {
      if app.ConflictsWith(req.TimeInfo) {
        return false
      }
  }
	return true
}

func (b *SolutionBoard) Assign(careGiver string, req IMRequest, idx int) {
  log.Trace("", "board.Assign", "Allocating request [%x] to %s", string(req.ID), careGiver)

  book := b.AllCareGivers[careGiver]
  book.Allocations = append(book.Allocations, req.TimeInfo)
  b.AllCareGivers[careGiver] = book

  b.CurrentSolution.Solution[idx].Status = "allocated"
  b.CurrentSolution.Solution[idx].CareGiver = careGiver
}

func (b *SolutionBoard) Unassign(careGiver string, req IMRequest, idx int) bool {
  for i, app := range b.AllCareGivers[careGiver].Allocations {
    if app.Equals(req.TimeInfo) {
      log.Trace("", "board.Unassign", "Unassigning request [%x] from %s", string(req.ID), careGiver)

      book := b.AllCareGivers[careGiver]
      book.Allocations = append(book.Allocations[:i], book.Allocations[i+1:]...)
      b.AllCareGivers[careGiver] = book

      b.CurrentSolution.Solution[idx].Status = "pending"
      b.CurrentSolution.Solution[idx].CareGiver = ""
      return true
    }
  }
  log.Trace("", "board.Unassign", "Unable to unassign request [%x] from %s!! Potential error emerging...", string(req.ID), careGiver)
  return false
}

func (req *IMRequest) Size() int {
	return len(req.CareGivers)
}


func (sol *SchedulingSolution) DumpToDB() {
  log.Trace("", "DumpToDB", "Should dump solution to DB")

  findService := &services.Service{}
  findService.UserID = "Scheduler subsystem"
  MonotonicMongoSession, err := mongo.CopyMonotonicSession("")
  if err != nil {
    log.Error(err, "Scheduler subsystem", "2) Couldn't copy monotonic session")
  }
  findService.MongoSession = MonotonicMongoSession

  for _, req := range sol.Solution {
    log.Trace("", "DumpToDB", "Permanent write allocation of [%s] to [%s] to DB", fmt.Sprintf("%x",string(req.ID)), req.CareGiver)
    err := requestService.AllocateRequest(findService,fmt.Sprintf("%x",string(req.ID)),req.CareGiver)
    if err != nil {
     log.CompletedErrorf(err, findService.UserID, "DumpToDB", "AllocateRequest")
    }
  }
}

func (sol *SchedulingSolution) Score() int {
  numPending := 0
  numTotal := len(sol.Solution)
  for _, req := range sol.Solution {
    if req.Status == "pending" {
        numPending += 1
    }
  }
  return (numTotal - numPending) * 100 / numTotal
}

func (sol *SchedulingSolution) Copy(s *SchedulingSolution) {
  newData := make([]IMRequest, 0)
  for _, req := range s.Solution {
      newData = append(newData, IMRequest{req.ID, nil, req.Status, req.CareGiver, req.TimeInfo})
  }
  sol.Solution = newData
}

func (sol *SchedulingSolution) Size() int {
  return len(sol.Solution)
}

var board *SolutionBoard

func killRecursion(dumpBest bool) {
  board.IsKilled = true
  if dumpBest {
    log.Trace("", "killRecursion", "Should dump best solution...")
    board.BestSolution.DumpToDB()
  }
}

func CreateSchedule() {
	log.Started("CreateSchedule", "Schedule creation started...")
  // First make sure previous scheduling process terminates:
  if board != nil {
    killRecursion(false)
  }

  // 1) Load existing solution from DB:
  // var findService hold a service pointer to perform queries on the service_requests db
  findService := &services.Service{}
  findService.UserID = "Scheduler subsystem"
  MonotonicMongoSession, err := mongo.CopyMonotonicSession("")
  if err != nil {
    log.Error(err, "Scheduler subsystem", "Couldn't copy monotonic session")
  }
  findService.MongoSession = MonotonicMongoSession

  log.Trace("", "CreateSchedule", "Fetching currentSchedule...")
  allRequests, err := requestService.FetchCurrentSchedule(findService)
  if err != nil {
    log.CompletedErrorf(err, findService.UserID, "CreateSchedule", "FetchRequest")
  }

  m := make(map[string]IMBookings)
  currentSolution := &SchedulingSolution{}
  bestSolution    := &SchedulingSolution{}
  for _, req := range allRequests {
    log.Trace("", "CreateSchedule", "Finding Care Givers for request [%x]", string(req.ID))
    // 2) As the solution is loading, fill in the possible care givers that could address each request:
    allCareGivers, err := userService.FetchPossibleCareGiversForRequest(findService, req)
    if err != nil {
      log.CompletedErrorf(err, findService.UserID, "CreateSchedule", "FetchPossibleCareGiversForRequest")
    }
    for _, cg := range allCareGivers {
      m[cg] = IMBookings{make([]Appointment, 0)}
    }
    currentSolution.Solution = append(currentSolution.Solution, IMRequest{req.ID, allCareGivers, "pending", "", Appointment{req.StartTime, req.Duration}})
    bestSolution.Solution    = append(   bestSolution.Solution, IMRequest{req.ID, nil, req.Status, req.CareGiver, Appointment{req.StartTime, req.Duration}})
  }

  board = &SolutionBoard{m, currentSolution, bestSolution, false}

  //Start Backtracking on currentSolution to find better solutions!
  log.Trace("", "CreateSchedule", "Starting recursion for depth[%d] with bestSolutionScore[%d] vs 0=%d?", len(board.CurrentSolution.Solution), bestSolution.Score(), currentSolution.Score())
  go play(0)
}

func play(req int) {
  log.Startedf("", "PLAY", "Running recursion level [%d]", req)

  if req < len(board.CurrentSolution.Solution) {
    lim := len(board.CurrentSolution.Solution[req].CareGivers)
    if lim > 0 {
      for cgIdx := 0; cgIdx < lim ; cgIdx++ {
          currentMove := board.CurrentSolution.Solution[req].CareGivers[cgIdx]
          if board.CanPlay(currentMove, board.CurrentSolution.Solution[req]) {
              board.Assign(currentMove, board.CurrentSolution.Solution[req], req)
              play(req + 1)
              board.Unassign(currentMove, board.CurrentSolution.Solution[req], req)
          } else {
              play(req + 1)
          }
      }
    } else {
      play(req + 1)
    }
  } else {
    log.Trace("", "PLAY", "Completed one branch! Assessing solution...")
    currentScore := board.CurrentSolution.Score()
    if currentScore > board.BestSolution.Score() {
          log.Trace("", "PLAY", "Backtracking found a better solution!! %d > %d", currentScore , board.BestSolution.Score())
          board.BestSolution.Copy(board.CurrentSolution)
          if currentScore > 99 {
            //This solution is too good. You may now stop the backtracking!
            killRecursion(true)
          } else {
            board.BestSolution.DumpToDB()
          }
    } else {
          log.Trace("", "PLAY", "Current solution is not better than existing best solution. CurrentScore[%d] <= %d", currentScore, board.BestSolution.Score())
    }
  }

  return
}
