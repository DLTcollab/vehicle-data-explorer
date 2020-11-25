$('.register_device_sub_form').on('click', function () {
    var data = getFormValue()
    axios({
        method: 'post',
        url: '/api/register_device',
        headers: {
            'Content-Type': 'application/json'
        },
        data
    }).then(res => {
        console.log(res.data)
        $('#message').text(res.data["message"])
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
