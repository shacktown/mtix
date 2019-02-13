/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');


const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', '..', 'basic-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

async function openWindow(contract, wName, loops) {
    try {
        var i;
        for (i=0; i<loops; i++) {
            var showNum = `show${getRandomInt(20)}`;
            var ticketNum = getRandomInt(5).toString();
            console.log(`show = ${showNum}  #tickets = ${ticketNum}`)
            var result = await contract.evaluateTransaction('get', showNum);
            console.log(`${wName} Transaction has been evaluated, result is: ${result.toString()}`);  
            var theShow = JSON.parse(result)
            console.log("Buying ", ticketNum, " tickets to see: ", theShow.Movie, "in ", theShow.Hall, " on ", theShow.DateTime)
           
            if (wName == "window1") {   
                result = await contract.submitTransaction("BuyTix", showNum, ticketNum, wName);
                console.log(`${wName} BuyTix Transaction has been submitted, result is: ${result.toString()}`);  
                result = await contract.submitTransaction("BuyConcession", "Regal1", "popcorn", ticketNum, theShow.DateTime);
                console.log(`${wName} BuyConcession Transaction has been submitted, result is: ${result.toString()}`);
                if ( getRandomInt(2) == 1 )  {
                    result = await contract.submitTransaction("ExchangeWaterSoda", "Regal1", theShow.DateTime);
                    console.log(`${wName} ExchangeWaterSoda Transaction has been submitted, result is: ${result.toString()}`);
                }
            }  
        }
    
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
 }

async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', 
                                    discovery: { enabled: false }, 
                                    eventHandlerOptions: {strategy: createTransactionEventHandler} }
                                    );

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('mtix');

 
       await openWindow(contract, "window1", 2);
       //openWindow(contract, "window2", 10);
       //openWindow(contract, "window3", 10);
       //await openWindow(contract, "window4", 10);
       console.log('Windows Closed.................');    

       // Disconnect from the gateway.
       await gateway.disconnect();
        

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

function createTransactionEventHandler(transactionId, network) {
    console.log("createTransactionEventHandler -------------------------------------------------");
	const channel = network.getChannel();
	const peers = channel.getPeersForOrg();
	const eventHubs = peers.map((peer) => channel.newChannelEventHub(peer.getName()));
	return new SampleTransactionEventHandler(transactionId, eventHubs);
}


function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
}
/************************************************************/
/**
 * Handler that listens for commit events for a specific transaction from a set of event hubs.
 * A new instance of this class should be created to handle each transaction as it maintains state
 * related to events for a given transaction.
 * @class
 */
class SampleTransactionEventHandler {
	/**
	 * Constructor.
	 * @param {String} transactionId Transaction ID for which events will be received.
	 * @param {ChannelEventHub[]} eventHubs Event hubs from which events will be received.
	 * @param {Object} [options] Additional configuration options.
	 * @param {Number} [options.commitTimeout] Time in seconds to wait for commit events to be reveived.
	 */
	constructor(transactionId, eventHubs, options) {
        console.log("constructor -------------------------------------------------");
		this.transactionId = transactionId;
		this.eventHubs = eventHubs;

		const defaultOptions = {
			commitTimeout: 120 // 2 minutes
		};
		this.options = Object.assign(defaultOptions, options);

		this.notificationPromise = new Promise((resolve, reject) => {
			this._resolveNotificationPromise = resolve;
			this._rejectNotificationPromise = reject;
		});
		this.eventCounts = {
			expected: this.eventHubs.length,
			received: 0
		};
	}

	/**
	 * Called to initiate listening for transaction events.
	 * @async
	 * @throws {Error} if not in a state where the handling strategy can be satified and the transaction should
	 * be aborted. For example, if insufficient event hubs could be connected.
	 */
	async startListening() {
        console.log("startListening -------------------------------------------------");
		if (this.eventHubs.length > 0) {
			this._setListenTimeout();
			await this._registerTxEventListeners();
		} else {
			// Assume success if no event hubs
			this._resolveNotificationPromise();
		}
	}

	/**
     * Wait until enough events have been received from the event hubs to satisfy the event handling strategy.
     * @async
	 * @throws {Error} if the transaction commit is not successful within the timeout period.
     */
	async waitForEvents() {
        console.log("waitForEvents -------------------------------------------------");
		await this.notificationPromise;
	}

	/**
     * Cancel listening for events.
     */
	cancelListening() {
		clearTimeout(this.timeoutHandler);
		this.eventHubs.forEach((eventHub) => {
			eventHub.unregisterTxEvent(this.transactionId);
			eventHub.disconnect();
		});
	}

	_setListenTimeout() {
		if (this.options.commitTimeout <= 0) {
			return;
		}

		this.timeoutHandler = setTimeout(() => {
			this._fail(new Error(`Timeout waiting for commit events for transaction ID ${this.transactionId}`));
		}, this.options.commitTimeout * 1000);
	}

	async _registerTxEventListeners() {
        console.log("_registerTxEventListeners -------------------------------------------------");
		const registrationOptions = {unregister: true, disconnect: true};

		const promises = this.eventHubs.map((eventHub) => {
			return new Promise((resolve) => {
				eventHub.registerTxEvent(
					this.transactionId,
					(txId, code) => this._onEvent(eventHub, txId, code),
					(err) => this._onError(eventHub, err),
					registrationOptions
				);
				eventHub.connect();
				resolve();
			});
		});

		await Promise.all(promises);
	}

	_onEvent(eventHub, txId, code) {
        console.log("_onEvent  txId=", txId, " code=", code);
		if (code !== 'VALID') {
			// Peer has rejected the transaction so stop listening with a failure
			const message = `Peer ${eventHub.getPeerAddr()} has rejected transaction ${txId} with code ${code}`;
			this._fail(new Error(message));
		} else {
			// --------------------------------------------------------------
			// Handle processing of successful transaction commit events here
			// --------------------------------------------------------------
			this._responseReceived();
		}
	}

	_onError(eventHub, err) { // eslint-disable-line no-unused-vars
		// ----------------------------------------------------------
		// Handle processing of event hub communication failures here
		// ----------------------------------------------------------
		this._responseReceived();
	}

	/**
	 * Simple event handling logic that is satisfied once all of the event hubs have either responded with valid
	 * events or disconnected.
	 */
	_responseReceived() {
		this.eventCounts.received++;
		if (this.eventCounts.received === this.eventCounts.expected) {
			this._success();
		}
	}

	_fail(error) {
		this.cancelListening();
		this._rejectNotificationPromise(error);
	}

	_success() {
		this.cancelListening();
		this._resolveNotificationPromise();
	}
}


/************************************************************/




main();
