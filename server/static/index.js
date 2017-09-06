let output = document.getElementById("output-text")

const message = (type, headerText, contentText) => {
  let element = document.createElement('div')
  element.classList.add('ui', type, 'message')

  let icon = document.createElement('i')
  icon.classList.add('close', 'icon')
  
  let header = document.createElement('div')
  header.classList.add('header')
  header.appendChild(document.createTextNode(headerText))
  
  let content = document.createElement('p')
  content.appendChild(document.createTextNode(contentText))

  element.appendChild(icon) 
  element.appendChild(header)
  element.appendChild(content)
  return element
}

const result = (headerText, contentText) => {
  document.getElementById('notifications').appendChild(message('success', headerText, contentText))
}

const error = (text) => {
  document.getElementById('notifications').appendChild(message('negative', headerText, contentText))
}

window.addEventListener('load', () => {
  init()
  addListeners()
})

const init = () => {
  //eth
  if (typeof web3 !== undefined) {
    window.web3 = new Web3(web3.currentProvider)
/*    switch (window.web3.version.network) {
      case "1":
        window.networkID = "Mainnet"
      case "3":
        window.networkID = "Ropsten"
      default:
        window.networkID = "Unknown"  
    }
    result.value = "Connected to " + window.networkID*/
  } else {
    window.networkID = undefined
  }

  //mint
  window.minturl = "http://localhost:12347"
//  $.get(minturl+"/keys", )
}

const addListeners = () => {
  $('.message .close').on('click', () => {
    $(this).closest('.message').transition('fade')
  })
  document.getElementById("setup-button").addEventListener("click", setupEvent)
  document.getElementById("deposit-button").addEventListener("click", depositEvent)
  document.getElementById("withdraw-button").addEventListener("click", withdrawEvent)
}

const setupEvent = () => {
  var flag = true
  let address = document.getElementById("contract-text")
  let keyname = document.getElementById("keyname-text")
  if (!web3.isAddress(address.value)) {
    document.getElementById("contract-input").classList.add("error")
    return
  } else {
    document.getElementById("contract-input").classList.remove("error")
  }
  $.getJSON("./keys/"+keyname.value, (key) => {
    if (key["success"] === false) {
      error("Error getting key from server", key["error"])
      flag = false
      return
    }
    window.keyname = keyname.value
  })
  $.getJSON("./ETGate.abi", (abi) => {
    let contract = window.web3.eth.contract(abi)
    window.instance = contract.at(address)
  })
  if (flag) result("Success to setup the key and contract", "Now you can deposit/withdraw your tokens")
}

const depositEvent = () => {
  if (window.instance == undefined) {
    result("Setup ETGate contract address first")
    return
  }

  let to = document.getElementById("to-text").value
  let value = Number(document.getElementById("value-text").value)
  let chain = "etgate-chain"
  let callback = (error, result) => {
      if (error) result("Error: " + error)
      else       result(result)
  }
  
  window.instance.depositEther(to, value, chain, {value: value}, callback)
}

const withdrawEvent = () => { 
  if (window.instance == undefined) {
    result("Setup ETGate contract address first")
    return
  }
}
