package retrier

import "github.com/taskie/zxc/model"

type SimpleRetrierState struct {
	LeftTimes int
}

type SimpleRetrier struct {
	state        *SimpleRetrierState
	config       *RetrierConfig
	globalConfig *model.GlobalConfig
}

func (retrier *SimpleRetrier) CheckIfRetriable() (bool, error) {
	return retrier.state.LeftTimes > 0, nil
}

func (retrier *SimpleRetrier) CheckIfRetriableNow() (bool, error) {
	ok, err := retrier.CheckIfRetriable()
	if err != nil || !ok {
		return ok, err
	}
	return true, nil
}

func (retrier *SimpleRetrier) Retry() error {
	retrier.state.LeftTimes--
	return nil
}
