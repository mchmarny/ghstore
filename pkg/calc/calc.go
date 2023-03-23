package calc

import (
	"strconv"

	"github.com/mchmarny/ghstore/pkg/data"
	"github.com/pkg/errors"
)

// Calculate performs the calculation based on the input map parameters
// and returns the result as a map.
func Calculate(m map[string]string) (map[string]string, error) {
	if m == nil {
		return nil, errors.New("input map is nil")
	}

	// init data store
	dataFile := m[DataKey]
	s, err := data.New(dataFile)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating data store: %s", dataFile)
	}
	defer s.Close()

	// get current counter value
	dataID := m[CounterKey]
	dataVal, err := s.Get(dataID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting data for: %s", dataID)
	}

	// parse input value
	dataValStr := m[ValueKey]
	newVal, err := strconv.ParseInt(dataValStr, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing value: %s", dataValStr)
	}

	opStr := m[ActionKey]
	op := ParseOperation(opStr)
	if op == OperationUndefined {
		return nil, errors.Errorf("invalid action: %s", opStr)
	}

	switch op {
	case OperationAdd:
		dataVal += newVal
	case OperationSub:
		dataVal -= newVal
	case OperationSet:
		dataVal = newVal
	case OperationGet:
	// do nothing
	default:
		return nil, errors.Errorf("invalid action: %s", op)
	}

	// save result
	if err := s.Upsert(dataID, dataVal); err != nil {
		return nil, errors.Wrapf(err, "error saving data: %s", dataID)
	}

	// get stored value
	storedVal, err := s.Get(dataID)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting data for: %s", dataID)
	}

	// validate result
	if storedVal != dataVal {
		return nil, errors.Errorf("error validating data: %s", dataID)
	}

	return map[string]string{
		ResultKey: strconv.FormatInt(storedVal, 10),
	}, nil
}
