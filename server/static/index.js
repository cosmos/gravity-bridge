let output = document.getElementById("output-text")

const etherAddress = "0x0000000000000000000000000000000000000000"

const message = (type, headerText, contentText) => {
  let element = document.createElement('div')
  element.classList.add('ui', type, 'message')

  let header = document.createElement('div')
  header.classList.add('header')
  header.appendChild(document.createTextNode(headerText))
  
  let content = document.createElement('p')
  content.appendChild(document.createTextNode(contentText))

  element.appendChild(header)
  element.appendChild(content)
  
  setTimeout(() => {
    element.remove()
  }, 5000)

  return element
}
const result = (headerText, contentText) => {
  document.getElementById('notifications').appendChild(message('success', headerText, contentText))
}

const error = (headerText, contentText) => {
  document.getElementById('notifications').appendChild(message('negative', headerText, contentText))
}

window.addEventListener('load', () => {
  init()
  addListeners()
})

const reloadBalance = () => {
  $.get("./query/account/"+window.key["name"], (resRaw) => {
    let res = JSON.parse(resRaw)
    let coin = res["result"].filter((x) => x["denom"] == encodeToken(etherAddress))[0]
    let balance = (coin || {"amount": 0})["amount"]

    let elem = document.getElementById("balance-text")
    elem.removeChild(elem.firstChild)
    elem.appendChild(document.createTextNode("Balance: " + web3.fromWei(balance, 'ether') + "ether"))
    
    setTimeout(() => {
      reloadBalance()
    }, 5000)
  })

}

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
    window.key = key
    reloadBalance()
  })
  $.getJSON("./ETGate.abi", (abi) => {
    let contract = window.web3.eth.contract(abi)
    window.instance = contract.at(address.value)
  })
  if (flag) result("Success to setup the key and contract", "Now you can deposit/withdraw your tokens")
}

const depositEvent = () => {
  if (window.instance == undefined) {
    result("Setup ETGate contract address first")
    return
  }

  let selection = document.getElementById("deposit-value-select")
  let selected = selection.options[selection.selectedIndex].value
  let valuetext = document.getElementById("deposit-value-text").value
  var value = 0
  if (selected != "wei") {
    value = Number(web3.toWei(valuetext, selected))
  } else {
    value = Number(valuetext)
  }
  let chain = "etgate-chain"
  let callback = (err, res) => {
      if (err) error("An error occured on deposit", err)
      else     result("Deposit successed", res)
  }
  
  window.instance.depositEther(("0x"+window.key["address"]).toLowerCase(), value, chain, {from: window.web3.defaultAccount, value: value}, callback)
}

const withdrawEvent = () => { 
  if (window.instance == undefined) {
    result("Setup ETGate contract address first")
    return
  }

  let selection = document.getElementById("withdraw-value-select")
  let selected = selection.options[selection.selectedIndex].value
  let valuetext = document.getElementById("withdraw-value-text").value
  var value = 0
  if (selected != "wei") {
    value = Number(web3.toWei(valuetext, selected))
  } else {
    value = Number(valuetext)
  }
  let chain = "etgate-chain"
  let to = document.getElementById("withdraw-to-text").value

  let data = {
    "name": window.key["name"],
    "passphrase": document.getElementById("withdraw-password-text").value,
    "to": to,
    "value": value,
    "token": etherAddress,
    "chainid": chain
  }

  let callback = (resRaw) => {
    let res = JSON.parse(resRaw)
    console.log(res)
//    withdrawEtherside(res["result"])
  }

  $.post("./withdraw", JSON.stringify(data), callback)

}

const withdrawEtherside = (data) => { 
  let callback = (err, res) => {
    if (err) {
      console.log(err)
      setTimeout(() => {
        withdrawEtherside(data)          
      }, 5000)
      return
    } 
    if (!res) {
      console.log("Not withdrawable, waiting...")
      setTimeout(() => {
        withdrawEtherside(data)
      }, 5000)
      return
    }

    let callback = (err, res) => {
      if (err) {
        console.log(err)
        error(err)
        setTimeout(() => {
          withdrawEtherside(data)
        }, 5000)
      } else {
        console.log(res)
        result(res)
      }
    }

    console.log(data)
    console.log(data["height"], 
                stringToBytes(data["iavlProofLeafHash"]), 
                data["iavlProofInnerHeight"],
                data["iavlProofInnerSize"],
                data["iavlProofInnerLeft"].map(stringToBytes),
                data["iavlProofInnerRight"].map(stringToBytes),
                stringToBytes(data["iavlProofRootHash"]),
                data["to"],
                data["value"],
                data["token"],
                data["chain"],
                data["seq"],
)

    window.instance.withdraw(data["height"], 
                             stringToBytes(data["iavlProofLeafHash"]), 
                             data["iavlProofInnerHeight"],
                             data["iavlProofInnerSize"],
                             data["iavlProofInnerLeft"].map(stringToBytes),
                             data["iavlProofInnerRight"].map(stringToBytes),
                             stringToBytes(data["iavlProofRootHash"]),
                             data["to"],
                             data["value"],
                             data["token"],
                             data["chain"],
                             data["seq"],
                             {},
                             callback)
  }
  window.instance.withdrawable(data["height"], data["chain"], data["token"], data["value"], {}, callback)
}

const encodeToken = (token) => {
  return token.slice(2).split('').map((x) => String.fromCharCode(x.charCodeAt()+32)).join('')
}

const decodeToken = (token) => {
  return "0x" + token.split('').map((x) => String.fromCharCode(x.charCodeAt()-32)).join('')
}

const stringToBytes = (str) => {
  return (str.match(/.{2}/g) || []).map((x) => "0x"+x)
}
