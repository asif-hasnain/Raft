package RPCs

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "raftAlgo.com/service/server/gRPC"
)

func (s *server) RequestAppendRPC(ctx context.Context, in *pb.RequestAppend) (*pb.ResponseAppend, error) {
	serverID := os.Getenv("CandidateID")
	log.Printf("Server %v : Received term : %v", serverID, in.GetTerm())
	log.Printf("Server %v : Received leaderId : %v", serverID, in.GetLeaderId())
	log.Printf("Server %v : Received prevLogIndex : %v", serverID, in.GetPrevLogIndex())
	log.Printf("Server %v : Received prevLogTerm : %v", serverID, in.GetPrevLogTerm())
	for i, entry := range in.GetEntries() {
		//TODO append log entry for worker
		log.Printf("Server %v : Received entry : %v at index %v", serverID, entry, i)
	}
	log.Printf("Server %v : Received leaderCommit : %v", serverID, in.GetLeaderCommit())
	//TODO commitIndex
	return &pb.ResponseAppend{Term: term, Success: false}, nil
}

//AppendRPC in leader
func AppendRPC(input string, address string) bool {
	// TODO go routine
	serverID := os.Getenv("CandidateID")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRPCServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//TODO change later
	// prevLogIndex = len(logs)
	// if lastLogIndex > 0 {
	// 	prevLogTerm = logs[lastLogIndex-1].term + 1
	// } else {
	// 	prevLogTerm = 1
	// }
	var prevLogIndex int64 = 2
	var prevLogTerm int64 = 2
	response, err := c.RequestAppendRPC(ctx, &pb.RequestAppend{Term: term, LeaderId: serverID, PrevLogIndex: prevLogIndex, PrevLogTerm: prevLogTerm,
		Entries: []*pb.RequestAppendLogEntry{{Command: input, Term: term}}, LeaderCommit: commitIndex})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Server %v : Response Received : %s", serverId, response.String())

	// TODO update leader data for each worker

	// TODO wait for >50% success response
	return true
}