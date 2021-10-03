# Functions:

- [`testMakeCheckpoint(address[] _validators, uint256[] _powers, uint256 _valsetNonce, bytes32 _gravityId)`](#Gravity-testMakeCheckpoint-address---uint256---uint256-bytes32-)

- [`testCheckValidatorSignatures(address[] _currentValidators, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, bytes32 _theHash, uint256 _powerThreshold)`](#Gravity-testCheckValidatorSignatures-address---uint256---uint8---bytes32---bytes32---bytes32-uint256-)

- [`lastBatchNonce(address _erc20Address)`](#Gravity-lastBatchNonce-address-)

- [`lastLogicCallNonce(bytes32 _invalidation_id)`](#Gravity-lastLogicCallNonce-bytes32-)

- [`updateValset(address[] _newValidators, uint256[] _newPowers, uint256 _newValsetNonce, address[] _currentValidators, uint256[] _currentPowers, uint256 _currentValsetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s)`](#Gravity-updateValset-address---uint256---uint256-address---uint256---uint256-uint8---bytes32---bytes32---)

- [`submitBatch(address[] _currentValidators, uint256[] _currentPowers, uint256 _currentValsetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256 _batchNonce, address _tokenContract, uint256 _batchTimeout)`](#Gravity-submitBatch-address---uint256---uint256-uint8---bytes32---bytes32---uint256---address---uint256---uint256-address-uint256-)

- [`submitLogicCall(address[] _currentValidators, uint256[] _currentPowers, uint256 _currentValsetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s, struct LogicCallArgs _args)`](#Gravity-submitLogicCall-address---uint256---uint256-uint8---bytes32---bytes32---struct-LogicCallArgs-)

- [`sendToCosmos(address _tokenContract, bytes32 _destination, uint256 _amount)`](#Gravity-sendToCosmos-address-bytes32-uint256-)

- [`deployERC20(string _cosmosDenom, string _name, string _symbol, uint8 _decimals)`](#Gravity-deployERC20-string-string-string-uint8-)

- [`constructor(bytes32 _gravityId, uint256 _powerThreshold, address[] _validators, uint256[] _powers)`](#Gravity-constructor-bytes32-uint256-address---uint256---)

# Events:

- [`TransactionBatchExecutedEvent(uint256 _batchNonce, address _token, uint256 _eventNonce)`](#Gravity-TransactionBatchExecutedEvent-uint256-address-uint256-)

- [`SendToCosmosEvent(address _tokenContract, address _sender, bytes32 _destination, uint256 _amount, uint256 _eventNonce)`](#Gravity-SendToCosmosEvent-address-address-bytes32-uint256-uint256-)

- [`ERC20DeployedEvent(string _cosmosDenom, address _tokenContract, string _name, string _symbol, uint8 _decimals, uint256 _eventNonce)`](#Gravity-ERC20DeployedEvent-string-address-string-string-uint8-uint256-)

- [`ValsetUpdatedEvent(uint256 _newValsetNonce, uint256 _eventNonce, address[] _validators, uint256[] _powers)`](#Gravity-ValsetUpdatedEvent-uint256-uint256-address---uint256---)

- [`LogicCallEvent(bytes32 _invalidationId, uint256 _invalidationNonce, bytes _returnData, uint256 _eventNonce)`](#Gravity-LogicCallEvent-bytes32-uint256-bytes-uint256-)

# Function `testMakeCheckpoint(address[] _validators, uint256[] _powers, uint256 _valsetNonce, bytes32 _gravityId)` {#Gravity-testMakeCheckpoint-address---uint256---uint256-bytes32-}

No description

# Function `testCheckValidatorSignatures(address[] _currentValidators, uint256[] _currentPowers, uint8[] _v, bytes32[] _r, bytes32[] _s, bytes32 _theHash, uint256 _powerThreshold)` {#Gravity-testCheckValidatorSignatures-address---uint256---uint8---bytes32---bytes32---bytes32-uint256-}

No description

# Function `lastBatchNonce(address _erc20Address) → uint256` {#Gravity-lastBatchNonce-address-}

No description

# Function `lastLogicCallNonce(bytes32 _invalidation_id) → uint256` {#Gravity-lastLogicCallNonce-bytes32-}

No description

# Function `updateValset(address[] _newValidators, uint256[] _newPowers, uint256 _newValsetNonce, address[] _currentValidators, uint256[] _currentPowers, uint256 _currentValsetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s)` {#Gravity-updateValset-address---uint256---uint256-address---uint256---uint256-uint8---bytes32---bytes32---}

No description

# Function `submitBatch(address[] _currentValidators, uint256[] _currentPowers, uint256 _currentValsetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s, uint256[] _amounts, address[] _destinations, uint256[] _fees, uint256 _batchNonce, address _tokenContract, uint256 _batchTimeout)` {#Gravity-submitBatch-address---uint256---uint256-uint8---bytes32---bytes32---uint256---address---uint256---uint256-address-uint256-}

No description

# Function `submitLogicCall(address[] _currentValidators, uint256[] _currentPowers, uint256 _currentValsetNonce, uint8[] _v, bytes32[] _r, bytes32[] _s, struct LogicCallArgs _args)` {#Gravity-submitLogicCall-address---uint256---uint256-uint8---bytes32---bytes32---struct-LogicCallArgs-}

No description

# Function `sendToCosmos(address _tokenContract, bytes32 _destination, uint256 _amount)` {#Gravity-sendToCosmos-address-bytes32-uint256-}

No description

# Function `deployERC20(string _cosmosDenom, string _name, string _symbol, uint8 _decimals)` {#Gravity-deployERC20-string-string-string-uint8-}

No description

# Function `constructor(bytes32 _gravityId, uint256 _powerThreshold, address[] _validators, uint256[] _powers)` {#Gravity-constructor-bytes32-uint256-address---uint256---}

No description

# Event `TransactionBatchExecutedEvent(uint256 _batchNonce, address _token, uint256 _eventNonce)` {#Gravity-TransactionBatchExecutedEvent-uint256-address-uint256-}

No description

# Event `SendToCosmosEvent(address _tokenContract, address _sender, bytes32 _destination, uint256 _amount, uint256 _eventNonce)` {#Gravity-SendToCosmosEvent-address-address-bytes32-uint256-uint256-}

No description

# Event `ERC20DeployedEvent(string _cosmosDenom, address _tokenContract, string _name, string _symbol, uint8 _decimals, uint256 _eventNonce)` {#Gravity-ERC20DeployedEvent-string-address-string-string-uint8-uint256-}

No description

# Event `ValsetUpdatedEvent(uint256 _newValsetNonce, uint256 _eventNonce, address[] _validators, uint256[] _powers)` {#Gravity-ValsetUpdatedEvent-uint256-uint256-address---uint256---}

No description

# Event `LogicCallEvent(bytes32 _invalidationId, uint256 _invalidationNonce, bytes _returnData, uint256 _eventNonce)` {#Gravity-LogicCallEvent-bytes32-uint256-bytes-uint256-}

No description
