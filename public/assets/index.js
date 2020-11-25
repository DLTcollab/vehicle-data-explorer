// hide the dashboard when loading at first time
document.getElementById("nav-item-dashboard").style.visibility = "hidden"

$('.mam_sub_form').on('click', function () {
    document.getElementById("spinner").style.visibility = "visible";
    document.getElementById("submit").disabled = true;
    var data = getFormValue()
    axios({
        method: 'post',
        url: '/api/sub_mam',
        headers: {
            'Content-Type': 'application/json'
        },
        data
    }).then(res => {
        console.log(res.data)
        if(res.data['status'] != 'success') {
            $('#mam_payload').append('<p>' + res.data['status'] + '</p>');
            return
        }

        for (var i = 0; i < res.data['data'].length; i++) {
            $('#mam_payload').append('<p>massage[' + i + "]" + "=" + res.data['data'][i] + '</p>');
        }

        document.getElementById("spinner").style.visibility = "hidden";
        document.getElementById("nav-item-dashboard").style.visibility = "visible"
        $('.json-msg').text(`success  ${new Date()}`)
    }).then(function () {
        document.getElementById("submit").disabled = false;
    });
})

function getFormValue() {
    var data = {}
    var inputs = $('#form input')
    for (let i = 0; i < inputs.length; i++) {
        data[$(inputs[i]).attr('name')] = $(inputs[i]).val()
    }
    return data
}
