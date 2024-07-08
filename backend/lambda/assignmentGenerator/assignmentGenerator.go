package assignmentGenerator

import (
	"fmt"
	"lambda/types"
	"lambda/util"
	"log"
	"math/rand"
)

func isAssignmentAllowed(giver string, receiver string, disallowedAssignments map[string][]string) bool {
	value, ok := disallowedAssignments[giver]
	if (ok && util.Contains(value, receiver)) || giver == receiver {
		log.Printf("Assignment from %s to %s is not allowed\n", giver, receiver)
		return false
	}
	log.Printf("Assignment from %s to %s is allowed\n", giver, receiver)
	return true
}

func isAssignmentComplete(assignments map[string]string) bool {
	for _, value := range assignments {
		if value == "" {
			log.Printf("Assignments are not complete: %s", assignments)
			return false
		}
	}
	log.Printf("Assignments are complete: %s", assignments)
	return true
}

func generateCandidateAssignments(currAssignment map[string]string, disallowedAssignments map[string][]string, giver string, unassignedReceivers []string) []map[string]string {
	var candidateAssignments []map[string]string

	for _, receiver := range unassignedReceivers {
		if isAssignmentAllowed(giver, receiver, disallowedAssignments) {
			candidate := util.Copy(currAssignment)
			candidate[giver] = receiver
			candidateAssignments = append(candidateAssignments, candidate)
		}
	}

	rand.Shuffle(len(candidateAssignments), func(i, j int) {
		candidateAssignments[i], candidateAssignments[j] = candidateAssignments[j], candidateAssignments[i]
	})

	log.Printf("Candidate assignments generated: %s", candidateAssignments)
	return candidateAssignments
}

func backtrackAssign(currAssignment map[string]string, disallowedAssignments map[string][]string, unassignedReceivers []string) bool {
	if isAssignmentComplete(currAssignment) {
		return true
	}

	giver := ""
	for key, value := range currAssignment {
		if value == "" {
			giver = key
			break
		}
	}

	candidateAssignments := generateCandidateAssignments(currAssignment, disallowedAssignments, giver, unassignedReceivers)
	for _, candidateAssignments := range candidateAssignments {
		unassignedReceivers = util.Remove(unassignedReceivers, candidateAssignments[giver])
		currAssignment[giver] = candidateAssignments[giver]
		if backtrackAssign(currAssignment, disallowedAssignments, unassignedReceivers) {
			return true
		}
		unassignedReceivers = append(unassignedReceivers, candidateAssignments[giver])
		currAssignment[giver] = ""
	}

	return false
}

func GenerateAssignments(group types.Group, restrictions map[string][]string) (map[string]string, error) {
	assignments := make(map[string]string)

	var unassignedReceivers []string
	for _, member := range group.GroupMembers {
		assignments[member.MemberId] = ""
		unassignedReceivers = append(unassignedReceivers, member.MemberId)
	}

	if backtrackAssign(assignments, restrictions, unassignedReceivers) {
		return assignments, nil
	} else {
		return assignments, fmt.Errorf("no valid assignments")
	}
}
