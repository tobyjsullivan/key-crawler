var bitcoin = require('bitcoinjs-lib');
var BigInteger = require('bigi');
var request = require('sync-request');

const BATCH_SIZE = 1000;
const QUEUER_HOSTNAME = process.env.QUEUER_HOSTNAME;
const QUEUER_PORT = process.env.QUEUER_PORT;
const PAIRS_ENDPOINT = `http://${QUEUER_HOSTNAME}:${QUEUER_PORT}/pairs`

var count = 0;

function enumerateKeys() {
    console.log("Entered enumerateKeys");
    var curPriv = BigInteger.ONE;

    var i = 0;
    while (i++ < 100) {
        var batch = generateBatch(curPriv);
        curPriv = curPriv.add(new BigInteger(''+BATCH_SIZE));

        var res = request('POST', PAIRS_ENDPOINT, {
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
    while (i++ < BATCH_SIZE) {
        var keyPair = new bitcoin.ECPair(curPriv);
        curPriv = curPriv.add(BigInteger.ONE);

        count++;
        batch.push({
            address: keyPair.getAddress(),
            'private-key': keyPair.toWIF()
        });

        if (count % 100 === 0) {
            console.log("count:", count);
        }
    }

    return batch;
}

enumerateKeys();
