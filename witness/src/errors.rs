use {web3, ethabi};
use std::io;

error_chain! {
    types {
        Error, ErrorKind, ResultExt, Result;
    }

    foreign_links {
        Io(io::Error);
        Ethabi(ethabi::Error);
    }

    errors {
        Web3(err: web3::Error) {
            description("web3 error"),
            display("{:?}", err),
        }
    }
}

