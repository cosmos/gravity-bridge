
e2e_setup:
	docker-compose down
	sudo rm -fr testdata

e2e_happy_path: e2e_setup
	./testnet/start.sh

e2e_validator_out: e2e_setup
	TEST_TYPE=VALIDATOR_OUT ./testnet/start.sh

e2e_batch_stress: e2e_setup
	TEST_TYPE=BATCH_STRESS ./testnet/start.sh

e2e_valset_stress: e2e_setup
	TEST_TYPE=VALSET_STRESS ./testnet/start.sh

e2e_v2_happy_path: e2e_setup
	TEST_TYPE=V2_HAPPY_PATH ./testnet/start.sh

e2e_arbitrary_logic: e2e_setup
	TEST_TYPE=ARBITRARY_LOGIC ./testnet/start.sh

e2e_orchestrator_keys: e2e_setup
	TEST_TYPE=ORCHESTRATOR_KEYS ./testnet/start.sh
