# CommandRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeciderTriggerType** | Pointer to [**DeciderTriggerType**](DeciderTriggerType.md) |  | [optional] 
**CommandWaitingType** | Pointer to [**CommandWaitingType**](CommandWaitingType.md) |  | [optional] 
**TimerCommands** | Pointer to [**[]TimerCommand**](TimerCommand.md) |  | [optional] 
**SignalCommands** | Pointer to [**[]SignalCommand**](SignalCommand.md) |  | [optional] 
**InterStateChannelCommands** | Pointer to [**[]InterStateChannelCommand**](InterStateChannelCommand.md) |  | [optional] 
**CommandCombinations** | Pointer to [**[]CommandCombination**](CommandCombination.md) |  | [optional] 

## Methods

### NewCommandRequest

`func NewCommandRequest() *CommandRequest`

NewCommandRequest instantiates a new CommandRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommandRequestWithDefaults

`func NewCommandRequestWithDefaults() *CommandRequest`

NewCommandRequestWithDefaults instantiates a new CommandRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDeciderTriggerType

`func (o *CommandRequest) GetDeciderTriggerType() DeciderTriggerType`

GetDeciderTriggerType returns the DeciderTriggerType field if non-nil, zero value otherwise.

### GetDeciderTriggerTypeOk

`func (o *CommandRequest) GetDeciderTriggerTypeOk() (*DeciderTriggerType, bool)`

GetDeciderTriggerTypeOk returns a tuple with the DeciderTriggerType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeciderTriggerType

`func (o *CommandRequest) SetDeciderTriggerType(v DeciderTriggerType)`

SetDeciderTriggerType sets DeciderTriggerType field to given value.

### HasDeciderTriggerType

`func (o *CommandRequest) HasDeciderTriggerType() bool`

HasDeciderTriggerType returns a boolean if a field has been set.

### GetCommandWaitingType

`func (o *CommandRequest) GetCommandWaitingType() CommandWaitingType`

GetCommandWaitingType returns the CommandWaitingType field if non-nil, zero value otherwise.

### GetCommandWaitingTypeOk

`func (o *CommandRequest) GetCommandWaitingTypeOk() (*CommandWaitingType, bool)`

GetCommandWaitingTypeOk returns a tuple with the CommandWaitingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandWaitingType

`func (o *CommandRequest) SetCommandWaitingType(v CommandWaitingType)`

SetCommandWaitingType sets CommandWaitingType field to given value.

### HasCommandWaitingType

`func (o *CommandRequest) HasCommandWaitingType() bool`

HasCommandWaitingType returns a boolean if a field has been set.

### GetTimerCommands

`func (o *CommandRequest) GetTimerCommands() []TimerCommand`

GetTimerCommands returns the TimerCommands field if non-nil, zero value otherwise.

### GetTimerCommandsOk

`func (o *CommandRequest) GetTimerCommandsOk() (*[]TimerCommand, bool)`

GetTimerCommandsOk returns a tuple with the TimerCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimerCommands

`func (o *CommandRequest) SetTimerCommands(v []TimerCommand)`

SetTimerCommands sets TimerCommands field to given value.

### HasTimerCommands

`func (o *CommandRequest) HasTimerCommands() bool`

HasTimerCommands returns a boolean if a field has been set.

### GetSignalCommands

`func (o *CommandRequest) GetSignalCommands() []SignalCommand`

GetSignalCommands returns the SignalCommands field if non-nil, zero value otherwise.

### GetSignalCommandsOk

`func (o *CommandRequest) GetSignalCommandsOk() (*[]SignalCommand, bool)`

GetSignalCommandsOk returns a tuple with the SignalCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalCommands

`func (o *CommandRequest) SetSignalCommands(v []SignalCommand)`

SetSignalCommands sets SignalCommands field to given value.

### HasSignalCommands

`func (o *CommandRequest) HasSignalCommands() bool`

HasSignalCommands returns a boolean if a field has been set.

### GetInterStateChannelCommands

`func (o *CommandRequest) GetInterStateChannelCommands() []InterStateChannelCommand`

GetInterStateChannelCommands returns the InterStateChannelCommands field if non-nil, zero value otherwise.

### GetInterStateChannelCommandsOk

`func (o *CommandRequest) GetInterStateChannelCommandsOk() (*[]InterStateChannelCommand, bool)`

GetInterStateChannelCommandsOk returns a tuple with the InterStateChannelCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterStateChannelCommands

`func (o *CommandRequest) SetInterStateChannelCommands(v []InterStateChannelCommand)`

SetInterStateChannelCommands sets InterStateChannelCommands field to given value.

### HasInterStateChannelCommands

`func (o *CommandRequest) HasInterStateChannelCommands() bool`

HasInterStateChannelCommands returns a boolean if a field has been set.

### GetCommandCombinations

`func (o *CommandRequest) GetCommandCombinations() []CommandCombination`

GetCommandCombinations returns the CommandCombinations field if non-nil, zero value otherwise.

### GetCommandCombinationsOk

`func (o *CommandRequest) GetCommandCombinationsOk() (*[]CommandCombination, bool)`

GetCommandCombinationsOk returns a tuple with the CommandCombinations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandCombinations

`func (o *CommandRequest) SetCommandCombinations(v []CommandCombination)`

SetCommandCombinations sets CommandCombinations field to given value.

### HasCommandCombinations

`func (o *CommandRequest) HasCommandCombinations() bool`

HasCommandCombinations returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


