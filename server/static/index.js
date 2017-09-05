let output = document.getElementById("output-text")

window.addEventListener('load', () => {
  init()
  addListeners()
})

const init = () => {
  //eth
  if (typeof web3 !== undefined) {
    window.web3 = new Web3(web3.currentProvider)
    switch (window.web3.version.network) {
      case "1":
        window.networkID = "Mainnet"
      case "3":
        window.networkID = "Ropsten"
      default:
        window.networkID = "Unknown"  
    }
    output.value = "Connected to " + window.networkID
  } else {
    window.networkID = undefined
  }

  $.getJSON("../contracts/ETGate.abi", (abi) => {
    let contract = window.web3.eth.contract(abi)
    let instance = contract.at(address)
  })

  //mint
  window.minturl = "http://localhost:12347"
  $.get(minturl+"/keys", )
}

const addListeners = () => {
  document.getElementById("setup-button").addEventListener("click", setupEvent)
  document.getElementById("deposit-button").addEventListener("click", depositEvent)
  document.getElementById("withdraw-button").addEventListener("click", withdrawEvent)
}

const setupEvent = () => {
  $.getJSON("../contracts/ETGate.abi", (abi) => {
    let address = document.getElementById("address-text").value
    let contract = window.web3.eth.contract(abi)
    window.instance = contract.at(address)
  })
}

const depositEvent = () => {
  if (window.instance == undefined) {
    output.value = "Setup ETGate contract address first"
    return
  }

  let to = document.getElementById("to-text").value
  let value = Number(document.getElementById("value-text").value)
  let chain = "etgate-chain"
  let callback = (error, result) => {
      if (error) output.value = "Error: " + error
      else       output.value = result
    }
  }
  window.instance.depositEther(to, value, chain, {value: value}, callback)
}

const withdrawEvent = () => { 
  if (window.instance == undefined) {
    output.value = "Setup ETGate contract address first"
    return
  }

  let tx = await window.client.build('app', {
     
  })
}
