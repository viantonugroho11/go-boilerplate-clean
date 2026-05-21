package sample

import (
	"context"

	entitysample "go-boilerplate-clean/internal/entity/sample"
	"go-boilerplate-clean/internal/usecase/sample/states"

	"gorm.io/gorm"
)

type (
	SampleAdder interface {
		Add(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error)
	}

	SampleGetter interface {
		Get(ctx context.Context, id string) (*entitysample.Sample, error)
	}

	SampleUpdater interface {
		Update(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error)
	}

	StateMachineSample interface {
		Do(ctx context.Context, tx *gorm.DB, update entitysample.Sample) (entitysample.Sample, error)
	}

	NewSampleStateMachine interface {
		NewStateMachine(ctx context.Context, current *entitysample.Sample) (states.ISampleStateMachine, error)
	}

	SamplePublisher interface {
		Publish(ctx context.Context, sample entitysample.Sample) error
	}

	TransactionManager interface {
		Begin(ctx context.Context) (*gorm.DB, error)
		Commit(ctx context.Context, tx *gorm.DB) error
		Rollback(ctx context.Context, tx *gorm.DB) error
	}

	// SampleService saves samples through the state machine orchestration.
	SampleService interface {
		Save(ctx context.Context, sample entitysample.Sample) (entitysample.Sample, error)
	}

)

type sampleSaver struct {
	stateMachine NewSampleStateMachine
	txManager    TransactionManager
	adder        SampleAdder
	getter       SampleGetter
	updater      SampleUpdater
	publisher    SamplePublisher
}

func NewSampleSaver(
	stateMachine NewSampleStateMachine,
	txManager TransactionManager,
	adder SampleAdder,
	getter SampleGetter,
	updater SampleUpdater,
	publisher SamplePublisher,
) SampleService {
	return &sampleSaver{
		stateMachine: stateMachine,
		txManager:    txManager,
		adder:        adder,
		getter:       getter,
		updater:      updater,
		publisher:    publisher,
	}
}

func (s *sampleSaver) Save(ctx context.Context, sample entitysample.Sample) (entitysample.Sample, error) {
	var (
		err     error
		updated entitysample.Sample
		current *entitysample.Sample
	)

	if sample.ID != "" {
		current, err = s.getter.Get(ctx, sample.ID)
		if err != nil {
			return entitysample.Sample{}, err
		}
	}

	storeFunc := s.updater.Update
	if current == nil {
		storeFunc = s.adder.Add
		current = &sample
	}

	stateMachine, err := s.stateMachine.NewStateMachine(ctx, current)
	if err != nil {
		return entitysample.Sample{}, err
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return entitysample.Sample{}, err
	}
	defer func() {
		if err != nil {
			_ = s.txManager.Rollback(ctx, tx)
		}
	}()

	updated, err = stateMachine.Do(ctx, tx, sample)
	if err != nil {
		return entitysample.Sample{}, err
	}

	updated, err = storeFunc(ctx, tx, updated)
	if err != nil {
		return entitysample.Sample{}, err
	}

	if err = s.txManager.Commit(ctx, tx); err != nil {
		return entitysample.Sample{}, err
	}

	if err := s.publisher.Publish(ctx, updated); err != nil {
		return entitysample.Sample{}, err
	}

	return updated, nil
}
