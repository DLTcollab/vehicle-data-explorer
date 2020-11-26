$('.btn_log_query').on('click', function () {
    var data = getDeviceID()
    axios({
        method: 'get',
        url: '/api/log/query',
        headers: {
            'Content-Type': 'application/json'
        },
        params: {
            deviceID: data["deviceID"]
        },
    }).then(res => {
        console.log(res.data)

        if (res.data['status'] === "success") {
            var logs = res.data['log']

            for (var i = 0; i < logs.length; i++) {
                var date = new Date(logs[i].timestamp * 1000);
                $('#logs').append('<p>massage[' + i + "]" + "=" + date + "," + logs[i].message + '</p>');
            }
        }
        $('#message').text(res.data["message"])
    }).then(function () {
        document.getElementById("submit").disabled = false;
    });
})

function getDeviceID() {
    var data = {}
    var inputs = $('#deviceID')
    data["deviceID"] = inputs.val()
    return data
}
