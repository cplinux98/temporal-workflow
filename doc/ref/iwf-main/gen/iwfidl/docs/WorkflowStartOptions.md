# WorkflowStartOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowIDReusePolicy** | Pointer to [**WorkflowIDReusePolicy**](WorkflowIDReusePolicy.md) |  | [optional] 
**CronSchedule** | Pointer to **string** |  | [optional] 
**WorkflowStartDelaySeconds** | Pointer to **int32** |  | [optional] 
**RetryPolicy** | Pointer to [**WorkflowRetryPolicy**](WorkflowRetryPolicy.md) |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**WorkflowConfigOverride** | Pointer to [**WorkflowConfig**](WorkflowConfig.md) |  | [optional] 
**IdReusePolicy** | Pointer to [**IDReusePolicy**](IDReusePolicy.md) |  | [optional] 
**UseMemoForDataAttributes** | Pointer to **bool** |  | [optional] 

## Methods

### NewWorkflowStartOptions

`func NewWorkflowStartOptions() *WorkflowStartOptions`

NewWorkflowStartOptions instantiates a new WorkflowStartOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStartOptionsWithDefaults

`func NewWorkflowStartOptionsWithDefaults() *WorkflowStartOptions`

NewWorkflowStartOptionsWithDefaults instantiates a new WorkflowStartOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowIDReusePolicy

`func (o *WorkflowStartOptions) GetWorkflowIDReusePolicy() WorkflowIDReusePolicy`

GetWorkflowIDReusePolicy returns the WorkflowIDReusePolicy field if non-nil, zero value otherwise.

### GetWorkflowIDReusePolicyOk

`func (o *WorkflowStartOptions) GetWorkflowIDReusePolicyOk() (*WorkflowIDReusePolicy, bool)`

GetWorkflowIDReusePolicyOk returns a tuple with the WorkflowIDReusePolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowIDReusePolicy

`func (o *WorkflowStartOptions) SetWorkflowIDReusePolicy(v WorkflowIDReusePolicy)`

SetWorkflowIDReusePolicy sets WorkflowIDReusePolicy field to given value.

### HasWorkflowIDReusePolicy

`func (o *WorkflowStartOptions) HasWorkflowIDReusePolicy() bool`

HasWorkflowIDReusePolicy returns a boolean if a field has been set.

### GetCronSchedule

`func (o *WorkflowStartOptions) GetCronSchedule() string`

GetCronSchedule returns the CronSchedule field if non-nil, zero value otherwise.

### GetCronScheduleOk

`func (o *WorkflowStartOptions) GetCronScheduleOk() (*string, bool)`

GetCronScheduleOk returns a tuple with the CronSchedule field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCronSchedule

`func (o *WorkflowStartOptions) SetCronSchedule(v string)`

SetCronSchedule sets CronSchedule field to given value.

### HasCronSchedule

`func (o *WorkflowStartOptions) HasCronSchedule() bool`

HasCronSchedule returns a boolean if a field has been set.

### GetWorkflowStartDelaySeconds

`func (o *WorkflowStartOptions) GetWorkflowStartDelaySeconds() int32`

GetWorkflowStartDelaySeconds returns the WorkflowStartDelaySeconds field if non-nil, zero value otherwise.

### GetWorkflowStartDelaySecondsOk

`func (o *WorkflowStartOptions) GetWorkflowStartDelaySecondsOk() (*int32, bool)`

GetWorkflowStartDelaySecondsOk returns a tuple with the WorkflowStartDelaySeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowStartDelaySeconds

`func (o *WorkflowStartOptions) SetWorkflowStartDelaySeconds(v int32)`

SetWorkflowStartDelaySeconds sets WorkflowStartDelaySeconds field to given value.

### HasWorkflowStartDelaySeconds

`func (o *WorkflowStartOptions) HasWorkflowStartDelaySeconds() bool`

HasWorkflowStartDelaySeconds returns a boolean if a field has been set.

### GetRetryPolicy

`func (o *WorkflowStartOptions) GetRetryPolicy() WorkflowRetryPolicy`

GetRetryPolicy returns the RetryPolicy field if non-nil, zero value otherwise.

### GetRetryPolicyOk

`func (o *WorkflowStartOptions) GetRetryPolicyOk() (*WorkflowRetryPolicy, bool)`

GetRetryPolicyOk returns a tuple with the RetryPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetryPolicy

`func (o *WorkflowStartOptions) SetRetryPolicy(v WorkflowRetryPolicy)`

SetRetryPolicy sets RetryPolicy field to given value.

### HasRetryPolicy

`func (o *WorkflowStartOptions) HasRetryPolicy() bool`

HasRetryPolicy returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowStartOptions) GetSearchAttributes() []SearchAttribute`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowStartOptions) GetSearchAttributesOk() (*[]SearchAttribute, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowStartOptions) SetSearchAttributes(v []SearchAttribute)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowStartOptions) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.

### GetWorkflowConfigOverride

`func (o *WorkflowStartOptions) GetWorkflowConfigOverride() WorkflowConfig`

GetWorkflowConfigOverride returns the WorkflowConfigOverride field if non-nil, zero value otherwise.

### GetWorkflowConfigOverrideOk

`func (o *WorkflowStartOptions) GetWorkflowConfigOverrideOk() (*WorkflowConfig, bool)`

GetWorkflowConfigOverrideOk returns a tuple with the WorkflowConfigOverride field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowConfigOverride

`func (o *WorkflowStartOptions) SetWorkflowConfigOverride(v WorkflowConfig)`

SetWorkflowConfigOverride sets WorkflowConfigOverride field to given value.

### HasWorkflowConfigOverride

`func (o *WorkflowStartOptions) HasWorkflowConfigOverride() bool`

HasWorkflowConfigOverride returns a boolean if a field has been set.

### GetIdReusePolicy

`func (o *WorkflowStartOptions) GetIdReusePolicy() IDReusePolicy`

GetIdReusePolicy returns the IdReusePolicy field if non-nil, zero value otherwise.

### GetIdReusePolicyOk

`func (o *WorkflowStartOptions) GetIdReusePolicyOk() (*IDReusePolicy, bool)`

GetIdReusePolicyOk returns a tuple with the IdReusePolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdReusePolicy

`func (o *WorkflowStartOptions) SetIdReusePolicy(v IDReusePolicy)`

SetIdReusePolicy sets IdReusePolicy field to given value.

### HasIdReusePolicy

`func (o *WorkflowStartOptions) HasIdReusePolicy() bool`

HasIdReusePolicy returns a boolean if a field has been set.

### GetUseMemoForDataAttributes

`func (o *WorkflowStartOptions) GetUseMemoForDataAttributes() bool`

GetUseMemoForDataAttributes returns the UseMemoForDataAttributes field if non-nil, zero value otherwise.

### GetUseMemoForDataAttributesOk

`func (o *WorkflowStartOptions) GetUseMemoForDataAttributesOk() (*bool, bool)`

GetUseMemoForDataAttributesOk returns a tuple with the UseMemoForDataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseMemoForDataAttributes

`func (o *WorkflowStartOptions) SetUseMemoForDataAttributes(v bool)`

SetUseMemoForDataAttributes sets UseMemoForDataAttributes field to given value.

### HasUseMemoForDataAttributes

`func (o *WorkflowStartOptions) HasUseMemoForDataAttributes() bool`

HasUseMemoForDataAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


