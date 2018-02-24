var bitcoin = require('bitcoinjs-lib');
var BigInteger = require('bigi');
var request = require('sync-request');

function enumerateKeys() {
    console.log("Entered enumerateKeys");
    var curPriv = BigInteger.ONE;

    var i = 0;
    while (i++ < 1000) {
        var batch = generateBatch(curPriv);
        curPriv.add(100);

        var res = request('POST', 'http://queuer:3000/pairs', {
            json: {
                pairs: batch
            }
        });
        console.log('Status:', res.statusCode);
    }
}

function generateBatch(start) {
    var curPriv = start;

    var batch = [];
    var i = 0;
    while (i++ < 100) {
        var keyPair = new bitcoin.ECPair(curPriv);
        curPriv = curPriv.add(BigInteger.ONE);
        // var curAddress = keyPair.getAddress();

        batch.push({
            address: keyPair.getAddress(),
            'private-key': keyPair.toWIF()
        });
    }

    return batch;
}

enumerateKeys();
