var bitcoin = require('bitcoinjs-lib');
var BigInteger = require('bigi');
var request = require('sync-request');
var AWS = require('aws-sdk');

AWS.config.update({region: 'us-east-1'});

const SQS_QUEUE_URL = process.env.SQS_QUEUE_URL;
const QUEUER_HOSTNAME = process.env.QUEUER_HOSTNAME;
const QUEUER_PORT = process.env.QUEUER_PORT;
const PAIRS_ENDPOINT = `http://${QUEUER_HOSTNAME}:${QUEUER_PORT}/pairs`;

function enumerateKeys() {
    console.log("Entered enumerateKeys");

    if (!SQS_QUEUE_URL) {
        throw Error("Must define SQS_QUEUE_URL");
    }

    const sqs = new AWS.SQS({apiVersion: '2012-11-05'});

    console.debug('[enumerateKeys] SQS client initialized.');

    setImmediate(processNextMessage, sqs);
}

function processNextMessage(sqs) {
    console.debug('[enumerateKeys] Requesting messages from SQS...');

    sqs.receiveMessage({
        QueueUrl: SQS_QUEUE_URL
    }, function(err, data) {
        console.debug('[enumerateKeys] receiveMessage callback entered.');

        if (err !== null) {
            console.error("[enumerateKeys] Error receiving messages:", err);
            return;
        }

        for (var msg of data.Messages) {
            const batchSpec = JSON.parse(msg.Body);

            var batch = generateBatch(batchSpec.start, batchSpec.size);
            var res = request('POST', PAIRS_ENDPOINT, {
                json: {
                    pairs: batch
                }
            });
            console.log('Status:', res.statusCode);

            if (res.statusCode !== 200) {
                console.error("[enumerateKeys] Error sending pairs:", res.statusMessage);
                continue;
            }

            sqs.deleteMessage({
                QueueUrl: SQS_QUEUE_URL,
                ReceiptHandle: msg.ReceiptHandle
            }, function(err) {
                if (err) {
                    console.error('[enumerateKeys] Error deleting message:', err);
                }
            });
        }

        setTimeout(processNextMessage, 1000, sqs);
    });

    console.debug('[enumerateKeys] Requested messages from SQS.');
}

function generateBatch(start, size) {
    var curPriv = new BigInteger('' + start);

    var batch = [];
    var i = 0;
    while (i++ < size) {
        var keyPair = new bitcoin.ECPair(curPriv);
        curPriv = curPriv.add(BigInteger.ONE);

        i++;
        batch.push({
            address: keyPair.getAddress(),
            'private-key': keyPair.toWIF()
        });

        if (i % 100 === 0) {
            console.log("count:", i, 'value:', curPriv);
        }
    }

    return batch;
}

enumerateKeys();
