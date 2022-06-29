package core

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
	"google.golang.org/protobuf/encoding/prototext"
)

type IndividualTestSuite struct {
	Name string

	// Names of all the tests, in the order they are tested.
	testNames []string

	testResults proto.TestSuiteResult
}

func NewIndividualTestSuite(suiteName string) *IndividualTestSuite {
	return &IndividualTestSuite{
		Name:        suiteName,
		testResults: newTestSuiteResult(),
	}
}

func (testSuite *IndividualTestSuite) TestCharacterStats(testName string, csr *proto.ComputeStatsRequest) {
	testSuite.testNames = append(testSuite.testNames, testName)

	result := ComputeStats(csr)
	finalStats := stats.FromFloatArray(result.RaidStats.Parties[0].Players[0].FinalStats)

	testSuite.testResults.CharacterStatsResults[testName] = &proto.CharacterStatsTestResult{
		FinalStats: finalStats[:],
	}
}

func (testSuite *IndividualTestSuite) TestStatWeights(testName string, swr *proto.StatWeightsRequest) {
	testSuite.testNames = append(testSuite.testNames, testName)

	result := StatWeights(swr)
	weights := stats.FromFloatArray(result.Dps.Weights)

	testSuite.testResults.StatWeightsResults[testName] = &proto.StatWeightsTestResult{
		Weights: weights[:],
	}
}

func (testSuite *IndividualTestSuite) TestDPS(testName string, rsr *proto.RaidSimRequest) {
	testSuite.testNames = append(testSuite.testNames, testName)

	result := RunRaidSim(rsr)
	if result.ErrorResult != "" {
		panic("simulation failed to run: " + result.ErrorResult)
	}
	testSuite.testResults.DpsResults[testName] = &proto.DpsTestResult{
		Dps:  result.RaidMetrics.Dps.Avg,
		Tps:  result.RaidMetrics.Parties[0].Players[0].Threat.Avg,
		Dtps: result.RaidMetrics.Parties[0].Players[0].Dtps.Avg,
	}
}

func (testSuite *IndividualTestSuite) Done(t *testing.T) {
	testSuite.writeToFile()
}

const tolerance = 0.00001

func (testSuite *IndividualTestSuite) writeToFile() {
	str := prototext.Format(&testSuite.testResults)
	// For some reason the formatter sometimes outputs 2 spaces instead of one.
	// Replace so we get consistent output.
	str = strings.ReplaceAll(str, "  ", " ")
	data := []byte(str)

	err := os.WriteFile(testSuite.Name+".results.tmp", data, 0644)
	if err != nil {
		panic(err)
	}
}

func (testSuite *IndividualTestSuite) readExpectedResults() proto.TestSuiteResult {
	data, err := os.ReadFile(testSuite.Name + ".results")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newTestSuiteResult()
		}

		panic(err)
	}

	results := &proto.TestSuiteResult{}
	if err = prototext.Unmarshal(data, results); err != nil {
		panic(err)
	}
	return *results
}

func newTestSuiteResult() proto.TestSuiteResult {
	return proto.TestSuiteResult{
		CharacterStatsResults: make(map[string]*proto.CharacterStatsTestResult),
		StatWeightsResults:    make(map[string]*proto.StatWeightsTestResult),
		DpsResults:            make(map[string]*proto.DpsTestResult),
	}
}

type TestGenerator interface {
	// The total number of tests that this generator can generate.
	NumTests() int

	// The name and API request for the test with the given index.
	GetTest(testIdx int) (string, *proto.ComputeStatsRequest, *proto.StatWeightsRequest, *proto.RaidSimRequest)
}

func RunTestSuite(t *testing.T, suiteName string, generator TestGenerator) {
	testSuite := NewIndividualTestSuite(suiteName)
	var currentTestName string

	defer func() {
		if p := recover(); p != nil {
			panic(fmt.Sprintf("Panic during test %s: %v", currentTestName, p))
		}
	}()

	expectedResults := testSuite.readExpectedResults()

	numTests := generator.NumTests()
	for i := 0; i < numTests; i++ {
		testName, csr, swr, rsr := generator.GetTest(i)
		if strings.Contains(testName, "Average") && testing.Short() {
			continue
		}
		currentTestName = testName

		t.Run(currentTestName, func(t *testing.T) {
			fullTestName := suiteName + "-" + testName
			if csr != nil {
				testSuite.TestCharacterStats(fullTestName, csr)
				if actualCharacterStats, ok := testSuite.testResults.CharacterStatsResults[fullTestName]; ok {
					actualStats := stats.FromFloatArray(actualCharacterStats.FinalStats)
					if expectedCharacterStats, ok := expectedResults.CharacterStatsResults[fullTestName]; ok {
						expectedStats := stats.FromFloatArray(expectedCharacterStats.FinalStats)
						if !actualStats.EqualsWithTolerance(expectedStats, tolerance) {
							t.Logf("Stats expected %v but was %v", expectedStats, actualStats)
							t.Fail()
						}
					} else {
						t.Logf("Unexpected test %s with stats: %v", fullTestName, actualStats)
						t.Fail()
					}
				} else if !ok {
					t.Logf("Missing Result for test %s", fullTestName)
					t.Fail()
				}
			} else if swr != nil {
				testSuite.TestStatWeights(fullTestName, swr)
				if actualStatWeights, ok := testSuite.testResults.StatWeightsResults[fullTestName]; ok {
					actualWeights := stats.FromFloatArray(actualStatWeights.Weights)
					if expectedStatWeights, ok := expectedResults.StatWeightsResults[fullTestName]; ok {
						expectedWeights := stats.FromFloatArray(expectedStatWeights.Weights)
						if !actualWeights.EqualsWithTolerance(expectedWeights, tolerance) {
							t.Logf("Weights expected %v but was %v", expectedWeights, actualWeights)
							t.Fail()
						}
					} else {
						t.Logf("Unexpected test %s with stat weights: %v", fullTestName, actualWeights)
						t.Fail()
					}
				} else if !ok {
					t.Logf("Missing Result for test %s", fullTestName)
					t.Fail()
				}
			} else if rsr != nil {
				testSuite.TestDPS(fullTestName, rsr)
				if actualDpsResult, ok := testSuite.testResults.DpsResults[fullTestName]; ok {
					if expectedDpsResult, ok := expectedResults.DpsResults[fullTestName]; ok {
						if actualDpsResult.Dps < expectedDpsResult.Dps-tolerance || actualDpsResult.Dps > expectedDpsResult.Dps+tolerance {
							t.Logf("DPS expected %0.03f but was %0.03f!.", expectedDpsResult.Dps, actualDpsResult.Dps)
							t.Fail()
						}
						if actualDpsResult.Tps < expectedDpsResult.Tps-tolerance || actualDpsResult.Tps > expectedDpsResult.Tps+tolerance {
							t.Logf("TPS expected %0.03f but was %0.03f!.", expectedDpsResult.Tps, actualDpsResult.Tps)
							t.Fail()
						}
						if actualDpsResult.Dtps < expectedDpsResult.Dtps-tolerance || actualDpsResult.Dtps > expectedDpsResult.Dtps+tolerance {
							t.Logf("DTPS expected %0.03f but was %0.03f!.", expectedDpsResult.Dtps, actualDpsResult.Dtps)
							t.Fail()
						}
					} else {
						t.Logf("Unexpected test %s with %0.03f DPS!", fullTestName, actualDpsResult.Dps)
						t.Fail()
					}
				} else if !ok {
					t.Logf("Missing Result for test %s", fullTestName)
					t.Fail()
				}
			} else {
				panic("No test request provided")
			}
		})
	}

	testSuite.Done(t)

	if t.Failed() {
		t.Log("One or more tests failed! If the changes are intentional, update the expected results with 'make test && make update-tests'. Otherwise go fix your bugs!")
	}
}
