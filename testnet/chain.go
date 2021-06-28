package main

import "fmt"

type Chain struct {
	DataDir    string
	ID         string
	Validators []*Validator
	Orchestrators []*Orchestrator
}

func (c *Chain) CreateAndInitializeValidators(count uint8) (err error){
	for i := uint8(0); i < count; i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		err = node.init()
		if err != nil {
			return
		}

		c.Validators = append(c.Validators, &node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}
	return
}

func (c *Chain) CreateAndInitializeValidatorsWithMnemonics(count uint8, mnemonics []string) (err error){
	for i := uint8(0); i < count; i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		err = node.init()
		if err != nil {
			return
		}

		c.Validators = append(c.Validators, &node)

		// create keys
		if err := node.createKeyFromMnemonic("val", mnemonics[i]); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}
	return
}

func (c *Chain) CreateAndInitializeOrchestrators(count uint8) (err error){
	for i := uint8(0); i < count; i++ {
		// create orchestrator
		orchestrator := c.createOrchestrator(i)

		// create keys
		mnemonic, info, err := createMemoryKey();
		if err != nil {
			return err
		}
		orchestrator.KeyInfo = *info
		orchestrator.Mnemonic = mnemonic

		c.Orchestrators = append(c.Orchestrators, &orchestrator)
	}
	return
}


func (c *Chain) CreateAndInitializeOrchestratorsWithMnemonics(count uint8, mnemonics []string) (err error){
	for i := uint8(0); i < count; i++ {
		// create orchestrator
		orchestrator := c.createOrchestrator(i)

		// create keys
		info, err := createMemoryKeyFromMnemonic(mnemonics[i]);
		if err != nil {
			return err
		}
		orchestrator.KeyInfo = *info
		orchestrator.Mnemonic = mnemonics[i]

		c.Orchestrators = append(c.Orchestrators, &orchestrator)
	}
	return
}

//func (c *Chain) RotateKeys() (err error) {
//	return
//}

func (c *Chain) createValidator(index uint8) (validator Validator) {
	validator = Validator{
		Chain:   c,
		Index:   index,
		Moniker: "gravity",
	}

	return
}

func (c *Chain) createOrchestrator(index uint8) (orchestrator Orchestrator) {
	orchestrator = Orchestrator{
		Chain:   c,
		Index:   index,
	}

	return
}

func (c *Chain) ConfigDir() string {
	return fmt.Sprintf("%s/%s", c.DataDir, c.ID)
}