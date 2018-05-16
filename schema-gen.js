var async = require("async")
var request = require("request")

var urlPrefix = "http://10.84.11.3:8081/analytics/table/"

var collections = [
    "StatTable.EndpointSecurityStats.eps.server",
    "StatTable.EndpointSecurityStats.eps.client",
    "StatTable.UveVirtualNetworkAgent.vn_stats"
    "StatTable.UveVMInterfaceAgent.if_stats",
    "StatTable.ComputeCpuState.cpu_info",
    "StatTable.NodeStatus.process_mem_cpu_usage"
    "StatTable.SandeshMessageStat.msg_info",
    "StatTable.VncApiStatsLog.api_stats",
    "StatTable.PeerStatsData.tx_update_stats",
    "StatTable.PeerStatsData.rx_update_stats",
    "StatTable.CollectorDbStats.table_info",
    "StatTable.QueryPerfInfo.query_stats",
    "StatTable.AnalyticsApiStats.api_stats",
    "StatTable.PRouterEntry.ifStats",
    "StatTable.CollectorDbStats.stats_info",
    "StatTable.CollectorDbStats.cql_stats"
];

var exceptList = ["COUNT_DISTINCT", "CLASS", "COUNT", "SUM", "MAX", "MIN",
    "PERCENTILES", "AVG"];

var collLen = collections.length;
var urlList = [];
for (var i = 0; i < collLen; i++) {
    urlList.push(urlPrefix + collections[i] + "/schema")
}
function getData(url, callback) {
    var options = {method: "GET", uri: url};
    request(options, function(error, response, data) {
        callback(error, JSON.parse(data));
    })
}

async.map(urlList, getData, function(error, data) {
    var count = data.length;
    var schemas = {}
    var exceptListCnt = exceptList.length;
    for (var i = 0; i < count; i++) {
        if (null == data[i]) {
            console.log("Getting NULL data for ", collections[i]);
            continue;
        }
        var cols = data[i].columns;
        if (0 == cols.length) {
            console.log("Columns are empty for ", collections[i]);
            continue;
        }
        schemas[collections[i]] = {}
        var colsLen = cols.length;
        for (var j = 0; j < colsLen; j++) {
            var name = cols[j].name;
            for (var k = 0; k < exceptListCnt; k++) {
               if (true == name.startsWith(exceptList[k])) {
                   break;
               }
            }
            if (k != exceptListCnt) {
                continue;
            }
            var datatype = cols[j].datatype;
            var index = cols[j].index;
            schemas[collections[i]][name] = datatype;
            if (true == index) {
                if (null == schemas[collections[i]]["tableindex"]) {
                    schemas[collections[i]]["tableindex"] = "";
                    schemas[collections[i]]["tableindex"] = name;
                } else {
                    schemas[collections[i]]["tableindex"] =
                        schemas[collections[i]]["tableindex"] + ";" + name;
                }
            }
        }
    }
    console.log(JSON.stringify(schemas, null, 4));
});
