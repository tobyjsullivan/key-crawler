var bitcoin = require('bitcoinjs-lib');
var BigInteger = require('bigi');
var request = require('sync-request');
var AWS = require('aws-sdk');

AWS.config.update({region: 'us-east-1'});

const BATCH_QUEUE_URL = process.env.BATCH_QUEUE_URL;
const QUEUER_HOSTNAME = process.env.QUEUER_HOSTNAME;
const QUEUER_PORT = process.env.QUEUER_PORT;
const PAIRS_ENDPOINT = `http://${QUEUER_HOSTNAME}:${QUEUER_PORT}/pairs`;

function enumerateKeys() {
    console.log("Entered enumerateKeys");

    if (!BATCH_QUEUE_URL) {
        throw Error("Must define BATCH_QUEUE_URL");
    }

    if (!QUEUER_HOSTNAME) {
        throw Error("Must define QUEUER_HOSTNAME");
    }

    if (!QUEUER_PORT) {
        throw Error("Must define QUEUER_PORT");
    }

    const sqs = new AWS.SQS({apiVersion: '2012-11-05'});

    console.debug('[enumerateKeys] SQS client initialized.');

    setImmediate(processNextMessage, sqs);
}

function processNextMessage(sqs) {
    console.debug('[processNextMessage] Requesting messages from SQS...');

    sqs.receiveMessage({
        QueueUrl: BATCH_QUEUE_URL
    }, function(err, data) {
        console.debug('[processNextMessage] receiveMessage callback entered.');

        if (err !== null) {
            console.error("[processNextMessage] Error receiving messages:", err);
            return;
        }

        console.debug("[processNextMessage] Received data:", data);
        console.debug("[processNextMessage] Received data.Messages:", data.Messages);

        for (let i = 0; i < data.Messages.length; i++) {
            const msg = data.Messages[i];
            const batchSpec = JSON.parse(msg.Body);

            const batchSize = 1000;
            for (let i = 0; i < batchSpec.size; i += batchSize) {
                var batch = generateBatch(batchSpec.start + i, batchSize);
                console.debug("[processNextMessage] PAIRS_ENDPOINT:", PAIRS_ENDPOINT);
                const payload = {
                    json: {
                        pairs: batch
                    }
                };
                var res = request('POST', PAIRS_ENDPOINT, payload);
                console.log('Status:', res.statusCode);

                if (res.statusCode !== 200) {
                    console.error("[processNextMessage] Error sending pairs:", res.statusMessage);
                    continue;
                }
            }

            sqs.deleteMessage({
                QueueUrl: BATCH_QUEUE_URL,
                ReceiptHandle: msg.ReceiptHandle
            }, function(err) {
                if (err) {
                    console.error('[processNextMessage] Error deleting message:', err);
                } else {
                    console.debug('[processNextMessage] Message deleted.');
                }
            });
        }

        setImmediate(processNextMessage, sqs);
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
