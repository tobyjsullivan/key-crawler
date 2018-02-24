var bitcoin = require('bitcoinjs-lib');
var BigInteger = require('bigi');
var request = require('sync-request');

function enumerateKeys() {
    console.log("Entered enumerateKeys");
    var curPriv = BigInteger.ONE;

    var i = 0;
    while (i++ < 1000) {
        var keyPair = new bitcoin.ECPair(curPriv);
        curPriv = curPriv.add(BigInteger.ONE);
        // var curAddress = keyPair.getAddress();

        console.log(`address=${keyPair.getAddress()}&private-key=${keyPair.toWIF()}`);
        var res = request('POST', 'http://localhost:3000/pairs', {
            json: {
                pairs: [{
                    address: keyPair.getAddress(),
                    'private-key': keyPair.toWIF()
                }]
            }
        });
        console.log('Status:', res.statusCode);
    }
}

enumerateKeys();
