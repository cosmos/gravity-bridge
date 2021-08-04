.DEFAULT_GOAL := e2e_slow_loris

e2e_build_images:
	@docker build -t gravity:prebuilt -f module/Dockerfile module/
	@docker build -t solidity:prebuilt -f solidity/Dockerfile solidity/
	@docker build -t orchestrator:prebuilt -f orchestrator/Dockerfile orchestrator/
	@docker build -t test-runner:prebuilt -f orchestrator/testnet.Dockerfile orchestrator/


e2e_slow_loris:
	@make -s e2e_happy_path
	@make -s e2e_v2_happy_path
	@make -s e2e_orchestrator_keys
	@make -s e2e_arbitrary_logic
	@make -s e2e_validator_out
	@make -s e2e_batch_stress
	@make -s e2e_valset_stress

e2e_clean_slate: e2e_build_images
	@docker rm --force \
		$(shell docker ps -qa --filter="name=contract_deployer") \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=gravity") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		$(shell docker ps -qa --filter="name=test_runner") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker wait \
		$(shell docker ps -qa --filter="name=contract_deployer") \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=gravity") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		$(shell docker ps -qa --filter="name=test_runner") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker network rm testnet 1>/dev/null 2>/dev/null || true
	@sudo rm -fr testdata
	@cd testnet && go test -c

e2e_batch_stress: e2e_clean_slate
	@testnet/testnet.test -test.run TestBatchStress -test.failfast -test.v || make -s fail

e2e_happy_path: e2e_clean_slate
	@testnet/testnet.test -test.run TestHappyPath -test.failfast -test.v || make -s fail

e2e_validator_out: e2e_clean_slate
	@testnet/testnet.test -test.run TestValidatorOut -test.failfast -test.v || make -s fail

e2e_valset_stress: e2e_clean_slate
	@testnet/testnet.test -test.run TestValsetStress -test.failfast -test.v || make -s fail

e2e_v2_happy_path: e2e_clean_slate
	@testnet/testnet.test -test.run TestV2HappyPath -test.failfast -test.v || make -s fail

e2e_arbitrary_logic: e2e_clean_slate
	@testnet/testnet.test -test.run TestArbitraryLogic -test.failfast -test.v || make -s fail

e2e_orchestrator_keys: e2e_clean_slate
	@testnet/testnet.test -test.run TestOrchestratorKeys -test.failfast -test.v || make -s fail

fail:
	@echo 'test failed; dumping container logs into ./testdata for review'
	@docker logs contract_deployer > testdata/contract_deployer.log 2>&1 || true
	@docker logs gravity0 > testdata/gravity0.log 2>&1 || true
	@docker logs gravity1 > testdata/gravity1.log 2>&1 || true
	@docker logs gravity2 > testdata/gravity2.log 2>&1 || true
	@docker logs gravity3 > testdata/gravity3.log 2>&1 || true
	@docker logs orchestrator0 > testdata/orchestrator0.log 2>&1 || true
	@docker logs orchestrator1 > testdata/orchestrator1.log 2>&1 || true
	@docker logs orchestrator2 > testdata/orchestrator2.log 2>&1 || true
	@docker logs orchestrator3 > testdata/orchestrator3.log 2>&1 || true
	@docker logs test_runner > testdata/test_runner.log 2>&1 || true
	@false
