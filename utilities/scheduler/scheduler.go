package scheduler

import (
  log "github.com/goinggo/tracelog"
  "github.com/goinggo/beego-mgo/utilities/availability"

  //Database access:
  "github.com/goinggo/beego-mgo/services"
	"github.com/goinggo/beego-mgo/services/userService"
	"github.com/goinggo/beego-mgo/services/requestService"
	"github.com/goinggo/beego-mgo/utilities/mongo"
  "gopkg.in/mgo.v2/bson"
)

type (
  IMCareGiver struct {
    ID         bson.ObjectId
    availability availability.Availability
  }

  IMRequest struct {
    ID         bson.ObjectId
    CareGivers []string //Care Givers identified by email!
    Status      string
    CareGiver   string
  }

  SchedulingSolution struct {
    Solution    []IMRequest
  }

  SolutionBoard struct {
    AllCareGivers map[string]bool //Merge of all Care Givers for all IMRequests! bool true means that Care Giver is available.
    CurrentSolution *SchedulingSolution
    BestSolution *SchedulingSolution
  }

)

func (req *IMRequest) CanPlay() bool {
	return len(req.CareGivers) > 0
}

func (req *IMRequest) Size() int {
	return len(req.CareGivers)
}


func (sol *SchedulingSolution) DumpToDB() {
  log.Trace("", "DumpToDB", "Should dump solution to DB")
  for _, req := range sol.Solution {
    log.Trace("", "DumpToDB", "Writing [%s] with status [%s] to DB", req.ID, req.Status)
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
  newData := make([]IMRequest, len(s.Solution))
  for _, req := range s.Solution {
      newData = append(newData, IMRequest{req.ID, nil, req.Status, req.CareGiver})
  }
  sol.Solution = newData
}

var board *SolutionBoard

func killRecursion(dumpBest bool) {
  for key, _ := range board.AllCareGivers {
    board.AllCareGivers[key] = false
  }

  if dumpBest {
    log.Trace("", "killRecursion", "Should dump best solution...")
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

  m := make(map[string]bool)
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
      m[cg] = true
    }
    currentSolution.Solution = append(currentSolution.Solution, IMRequest{req.ID, allCareGivers, "pending", ""})
    bestSolution.Solution    = append(   bestSolution.Solution, IMRequest{req.ID, nil, req.Status, req.CareGiver})
  }

  board = &SolutionBoard{m, currentSolution, bestSolution}

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
          //If you can play this move, play it:
          if board.AllCareGivers[currentMove] {

            //log.Trace("", "PLAY", "Allocating request [%x] to %s", string(board.CurrentSolution.Solution[req].ID), board.CurrentSolution.Solution[req].CareGivers[cgIdx])
            board.AllCareGivers[currentMove] = false
            board.CurrentSolution.Solution[req].Status = "allocated"
            board.CurrentSolution.Solution[req].CareGiver = board.CurrentSolution.Solution[req].CareGivers[cgIdx]

            play(req + 1)

            //Undo:
            board.AllCareGivers[currentMove] = true
            board.CurrentSolution.Solution[req].Status = "pending"
            board.CurrentSolution.Solution[req].CareGiver = ""
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
          }
    } else {
          log.Trace("", "PLAY", "Current solution is not better than existing best solution. CurrentScore[%d]", currentScore)
    }
  }

  return
}
