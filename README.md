# DynamoDB and PostgreSQL Benchmarks

Most useful columns to look at first sight:

* Gophers: number of concurrent request at the same time;
* Time Per Action: time per action;

In both the write and read tests, PostgreSQL is having better performance in general.

One thing to mind is that in the go benchmark programs, the AWS DynamoDB SDK is using HTTP (and JSON, maybe) to save and retrieve data, while the PostgreSQL is communicating with its own server over their own protocol (at least not HTTP[S] which I am sure). But this should be justified is because our server is going to written in Go and the DynamoDB SDK is from Amazon.

# Write

PostgreSQL has better write performance over DynamoDB (WCU: 1000).

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
PostgreSQL|8|15000000|1h51m59.261364771s|3.576059ms|14h54m0.888903876s
DynamoDB (WCU: 1000)|8|20000000|8h33m38.893351425s|6.158578ms|34h12m51.567642963s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
PostgreSQL|1|100000|2m12.735682795s|1.318288ms|2m11.828872449s
PostgreSQL|4|200000|1m49.217753069s|2.175379ms|7m15.075873342s
PostgreSQL|16|400000|2m24.950670158s|5.790638ms|38m36.255239229s
PostgreSQL|64|400000|2m38.64933233s|25.366757ms|2h49m6.702871642s
PostgreSQL|128|400000|2m42.108330581s|51.663684ms|5h44m25.473883248s
PostgreSQL|256|400000|2m38.41617564s|100.1726ms|11h7m49.040073977s
PostgreSQL|512|400000|2m39.156508093s|198.30337ms|22h2m1.348170585s
DynamoDB (WCU: 1000)|1|40000|4m33.860506028s|6.841327ms|4m33.65308054s
DynamoDB (WCU: 1000)|4|40000|1m5.215556218s|6.515916ms|4m20.636652622s
DynamoDB (WCU: 1000)|16|300000|2m19.875626219s|7.452623ms|37m15.786926217s
DynamoDB (WCU: 1000)|64|300000|2m58.181090652s|37.91286ms|3h9m33.858236529s
DynamoDB (WCU: 1000)|128|300000|2m34.260546971s|63.93929ms|5h19m41.787248308s
DynamoDB (WCU: 1000)|256|300000|4m26.385739836s||223.045787ms18h35m13.736355669s

# Read

PostgreSQL is having better performance over DynamoDB.

But while receiving 256 concurrent request, PostgreSQL starts returning error: `pq: remaining connection slots are reserved for non-replication superuser connections`. DynamoDB also returned error like `connect: cannot assign requested address
RequestError: send request failed` for 256 concurrent requests.

While for DynamoDB, query benchmark returned unstable result over the tests with concurrent requests number larger than 64 (marked with * and each has detail result at the bottom).

PostgreSQL Database Size: 21866690

DynamoDB Database Size:   20000000

(RCU represents: [Read Capacity Unit]())

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
PostgreSQL|1|400000|3m34.036918517s|527.178µs|3m30.871550027s
PostgreSQL|4|400000|1m12.09766096s|713.366µs|4m45.34660808s
PostgreSQL|16|800000|2m4.208954425s|2.47646ms|33m1.168660067s
PostgreSQL|64|800000|2m10.728693997s|10.444029ms|2h19m15.223794865s
PostgreSQL|128|600000|1m47.880797788s|22.916718ms|3h49m10.031389567s
PostgreSQL(Err)|256|500000|1m48.850811949s|55.304314ms|7h40m52.157350963s
DynamoDB (RCU: 1000)|1|40000|3m11.882240877s|4.791944ms|3m11.677760542s
DynamoDB (RCU: 1000)|4|100000|1m58.873632475s|4.749514ms|7m54.951446257s
DynamoDB (RCU: 1000)|16|400000|3m32.197960362s|8.478758ms|56m31.50359683s
DynamoDB (RCU: 1000)*|64|250000|1m18.201855788s|25.571496ms|1h22m19.264477631s
DynamoDB (RCU: 1000)*|128|260000|2m0.268183173s|68.446384ms|4h12m4.269993534s
DynamoDB (RCU: 1000)*|256|450000|6m10.728970646s|207.285081ms|25h44m36.73005557s
DynamoDB (RCU: 5000)|1|40000|2m16.74213163s|3.413038ms|2m16.521524267s
DynamoDB (RCU: 5000)|4|100000|1m42.025945416s|4.075698ms|6m47.56987479s
DynamoDB (RCU: 5000)|16|400000|2m31.908523721s|6.067671ms|40m27.068745639s
DynamoDB (RCU: 5000)*|64|400000|3m6.844018263s|29.731645ms|3h18m12.658443519s
DynamoDB (RCU: 5000)*|128|500000|5m2.89205501s|73.755437ms|10h32m57.733658743s
DynamoDB (RCU: 5000)*|256|466666|4m10.980653186s|131.074171ms|17h19m13.151302908s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
DynamoDB (RCU: 1000)|64|400000|1m25.373076402s|13.554102ms|1h30m21.640854546s
DynamoDB (RCU: 1000)|64|400000|1m54.537359725s|18.207428ms|2h1m22.971520966s
DynamoDB (RCU: 1000)|64|100000|53.227440049s|32.478477ms|54m7.847754072s
DynamoDB (RCU: 1000)|64|100000|59.669546979s|38.045977ms|1h3m24.597780943s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
DynamoDB (RCU: 1000)|128|100000|59.276534133s|74.915274ms|2h4m51.527427023s
DynamoDB (RCU: 1000)|128|200000|3m13.146200035s|122.591917ms|6h48m38.383432873s
DynamoDB (RCU: 1000)|128|200000|1m4.407081845s|40.312814ms|2h14m22.562836411s
DynamoDB (RCU: 1000)|128|200000|1m47.318676591s|67.745684ms|3h45m49.136815463s
DynamoDB (RCU: 1000)|128|600000|2m57.192423262s|36.666232ms|6h6m39.739455904s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
DynamoDB (RCU: 1000)|256|800000|9m55.959232044s|185.822585ms|41h17m38.068065537s
DynamoDB (RCU: 1000)|256|400000|5m50.649290392s|219.034951ms|24h20m13.980434483s
DynamoDB (RCU: 1000)|256|400000|6m34.159067526s|247.891568ms|27h32m36.627559989s
DynamoDB (RCU: 1000)|256|200000|2m22.148292624s|176.39122ms|9h47m58.244162273s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
DynamoDB (RCU: 5000)|64|400000|2m6.433985059s|20.066885ms|2h13m46.754002357s
DynamoDB (RCU: 5000)|64|400000|3m42.90090105s|35.408459ms|3h56m3.383656876s
DynamoDB (RCU: 5000)|64|400000|2m24.692404994s|23.050429ms|2h33m40.171761834s
DynamoDB (RCU: 5000)|64|400000|4m13.348781951s|40.40081ms|4h29m20.324353009s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
DynamoDB (RCU: 5000)|128|400000|1m37.407576096s|30.266023ms|3h21m46.409364774s
DynamoDB (RCU: 5000)|128|600000|5m18.301101638s|66.618784ms|11h6m11.270799489s
DynamoDB (RCU: 5000)|128|600000|8m9.533071717s|102.892388ms|17h8m55.433207371s
DynamoDB (RCU: 5000)|128|400000|5m6.326470592s|95.244553ms|10h34m57.82126334s

Type|Gophers|Request Count|Total Time|Time Per Action|Total Duration
----|----|----|----|
DynamoDB (RCU: 5000)|256|800000|7m56.815498277s|148.852239ms|33h4m41.79195969s
DynamoDB (RCU: 5000)|256|400000|4m17.034918153s|160.130754ms|17h47m32.30182556s
DynamoDB (RCU: 5000)|256|400000|2m32.608083526s|94.025159ms|10h26m50.063731445s
DynamoDB (RCU: 5000)|256|400000|3m40.249139529s|136.651411ms|15h11m0.564716366s
DynamoDB (RCU: 5000)|256|400000|3m57.214627082s|147.321619ms|16h22m8.64761368s
DynamoDB (RCU: 5000)|256|400000|2m41.96165255s|99.463844ms|11h3m5.537970708s
