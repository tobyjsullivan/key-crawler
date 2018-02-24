var bitcoin = require('bitcoinjs-lib');
var BigInteger = require('bigi');


function enumerateKeys() {
	var curPriv = BigInteger.ONE;

	var i = 0;
	while (true) {
		var keyPair = new bitcoin.ECPair(curPriv);
		curPriv = curPriv.add(BigInteger.ONE);
		var curAddress = keyPair.getAddress();

		console.log(keyPair.getAddress()+","+keyPair.toWIF())
	}
}

enumerateKeys();
