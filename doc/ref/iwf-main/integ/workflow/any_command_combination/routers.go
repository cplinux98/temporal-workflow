package anycommandcombination

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
	"time"
)

const (
	WorkflowType     = "any_command_combination"
	State1           = "S1"
	State2           = "S2"
	TimerId1         = "test-timer-1"
	SignalNameAndId1 = "test-signal-name1"
	SignalNameAndId2 = "test-signal-name2"
	SignalNameAndId3 = "test-signal-name3"
)

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
	//we want to confirm that the interpreter workflow activity will fail when the commandId is empty with ANY_COMMAND_COMBINATION_COMPLETED
	hasS1RetriedForInvalidCommandId bool
	hasS2RetriedForInvalidCommandId bool
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory:                   make(map[string]int64),
		invokeData:                      make(map[string]interface{}),
		hasS1RetriedForInvalidCommandId: false,
		hasS2RetriedForInvalidCommandId: false,
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	invalidTimerCommands := []iwfidl.TimerCommand{
		{
			CommandId:                  "",
			FiringUnixTimestampSeconds: iwfidl.PtrInt64(time.Now().Unix() + 86400*365), // one year later
		},
	}
	validTimerCommands := []iwfidl.TimerCommand{
		{
			CommandId:                  TimerId1,
			FiringUnixTimestampSeconds: iwfidl.PtrInt64(time.Now().Unix() + 86400*365), // one year later
		},
	}
	invalidSignalCommands := []iwfidl.SignalCommand{
		{
			CommandId:         "",
			SignalChannelName: SignalNameAndId1,
		},
		{
			CommandId:         SignalNameAndId2,
			SignalChannelName: SignalNameAndId2,
		},
	}
	validSignalCommands := []iwfidl.SignalCommand{
		{
			CommandId:         SignalNameAndId1,
			SignalChannelName: SignalNameAndId1,
		},
		{
			CommandId:         SignalNameAndId1,
			SignalChannelName: SignalNameAndId1,
		},
		{
			CommandId:         SignalNameAndId2,
			SignalChannelName: SignalNameAndId2,
		},
		{
			CommandId:         SignalNameAndId3,
			SignalChannelName: SignalNameAndId3,
		},
	}

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			if h.hasS1RetriedForInvalidCommandId {
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     validSignalCommands,
						TimerCommands:      validTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
						CommandCombinations: []iwfidl.CommandCombination{
							{
								CommandIds: []string{
									TimerId1, SignalNameAndId1, SignalNameAndId1, // wait for two SignalNameAndId1
								},
							},
							{
								CommandIds: []string{
									TimerId1, SignalNameAndId1, SignalNameAndId2,
								},
							},
						},
					},
				}

				c.JSON(http.StatusOK, startResp)
			} else {
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     validSignalCommands,
						TimerCommands:      invalidTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
					},
				}
				h.hasS1RetriedForInvalidCommandId = true
				c.JSON(http.StatusOK, startResp)
			}
			return
		}
		if req.GetWorkflowStateId() == State2 {
			if h.hasS2RetriedForInvalidCommandId {
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     validSignalCommands,
						TimerCommands:      validTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
						CommandCombinations: []iwfidl.CommandCombination{
							{
								CommandIds: []string{
									SignalNameAndId2, SignalNameAndId1,
								},
							},
							{
								CommandIds: []string{
									TimerId1, SignalNameAndId1, SignalNameAndId2,
								},
							},
						},
					},
				}

				c.JSON(http.StatusOK, startResp)
			} else {
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     invalidSignalCommands,
						TimerCommands:      validTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
					},
				}
				h.hasS2RetriedForInvalidCommandId = true
				c.JSON(http.StatusOK, startResp)
			}
			return
		}
	}

	panic("invalid workflow type")
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		if req.GetWorkflowStateId() == State1 {
			h.invokeData["s1_commandResults"] = req.GetCommandResults()

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State2,
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			h.invokeData["s2_commandResults"] = req.GetCommandResults()

			// go to complete
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: service.GracefulCompletingWorkflowStateId,
						},
					},
				},
			})
			return
		}
	}

	panic("invalid workflow type")
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
