var obd2_datasets;

var engineChartData = [3276.75, 6553.5, 9830.25, 16383.75]
var engineChartValue = 1000

var EngineConfig = {
    type: 'gauge',
    data: {
        //labels: ['Success', 'Warning', 'Warning', 'Error'],
        datasets: [{
            data: engineChartData,
            value: engineChartValue,
            backgroundColor: ['green', 'yellow', 'orange', 'red'],
            borderWidth: 2
        }]
    },
    options: {
        responsive: true,
        title: {
            display: true,
            text: 'Engine Speed(RPM)'
        },
        layout: {
            padding: {
                bottom: 30
            }
        },
        needle: {
            // Needle circle radius as the percentage of the chart area width
            radiusPercentage: 2,
            // Needle width as the percentage of the chart area width
            widthPercentage: 3.2,
            // Needle length as the percentage of the interval between inner radius (0%) and outer radius (100%) of the arc
            lengthPercentage: 80,
            // The color of the needle
            color: 'rgba(0, 0, 0, 1)'
        },
        valueLabel: {
            formatter: Math.round
        }
    }
};

var VehicleChartData = [51, 102, 153, 255]
var VehicleChartValue = 40

var VehicleConfig = {
    type: 'gauge',
    data: {
        //labels: ['Success', 'Warning', 'Warning', 'Error'],
        datasets: [{
            data: VehicleChartData,
            value: VehicleChartValue,
            backgroundColor: ['green', 'yellow', 'orange', 'red'],
            borderWidth: 2
        }]
    },
    options: {
        responsive: true,
        title: {
            display: true,
            text: '	Vehicle speed(Km/h)'
        },
        layout: {
            padding: {
                bottom: 30
            }
        },
        needle: {
            // Needle circle radius as the percentage of the chart area width
            radiusPercentage: 2,
            // Needle width as the percentage of the chart area width
            widthPercentage: 3.2,
            // Needle length as the percentage of the interval between inner radius (0%) and outer radius (100%) of the arc
            lengthPercentage: 80,
            // The color of the needle
            color: 'rgba(0, 0, 0, 1)'
        },
        valueLabel: {
            formatter: Math.round
        }
    }
};

var barChartData = {
    type: 'bar',
    data: {
        labels: ['Engine coolant temperature', 'Intake air temperature', 'Ambient air temperature', 'Engine oil temperature'],
        datasets: [{
            label: 'Temperature(°C)',
            data: [12, 19, 3, 3],
            backgroundColor: [
                'rgba(255, 99, 132, 0.2)',
                'rgba(54, 162, 235, 0.2)',
                'rgba(255, 206, 86, 0.2)',
                'rgba(75, 192, 192, 0.2)',
            ],
            borderColor: [
                'rgba(255, 99, 132, 1)',
                'rgba(54, 162, 235, 1)',
                'rgba(255, 206, 86, 1)',
                'rgba(75, 192, 192, 1)',
            ],
            borderWidth: 1
        }]
    }
}

var barChartOptions = {
    scales: {
        yAxes: [{
            ticks: {
                beginAtZero: true
            }
        }]
    }
}

var lineChartData = {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: "Time series data",
            data: [],
            lineTension: 0,
            backgroundColor: 'transparent',
            borderColor: '#007bff',
            borderWidth: 2,
            pointBackgroundColor: '#007bff',
        }]
    },
    options: {
        scales: {
            xAxes: [{
                type: 'time',
                distribution: 'series'
            }]
        }
    }
}

function get_dashboard_data() {
    return axios({
        method: 'get',
        url: '/api/dashboard_data',
        headers: {
            'Content-Type': 'application/json'
        },
    }).then(res => {
        console.log(res.data)
        var datasets = res.data['data']
        $('.json-msg').text(`success  ${new Date()}`)
        datasets.sort((a, b) => a.timestamp - b.timestamp);
        return datasets;
    })
}

var engineGaugeCtx = document.getElementById('Engine_speed_chart').getContext('2d');
window.engineGauge = new Chart(engineGaugeCtx, EngineConfig);

var VehicleGaugeCtx = document.getElementById('Vehicle_speed_chart').getContext('2d');
window.VehicleGauge = new Chart(VehicleGaugeCtx, VehicleConfig);

var barChartCtx = document.getElementById('temperature_chart');
window.barChart = new Chart(barChartCtx, barChartData, barChartOptions);

var lineChartCtx = document.getElementById('lineChart');
window.lineChart = new Chart(lineChartCtx, lineChartData);

function update_dashboard_components(datasets) {
    var json_datasets = datasets

    // fetch the latest data
    var bind_data = json_datasets[json_datasets.length - 1].data

    if ($('#dropdown-menu').children().length == 0) {
        update_drop_down_labels(bind_data)
    }

    document.getElementById('text_vin').textContent = bind_data["vin"]
    document.getElementById('text_engine_load').textContent = bind_data['engine_load']
    document.getElementById('text_mass_air_flow').textContent = bind_data['mass_air_flow']
    document.getElementById('text_fuel_tank_level_input').textContent = bind_data['fuel_tank_level_input']
    document.getElementById('text_control_module_voltage').textContent = bind_data['control_module_voltage']
    document.getElementById('text_throttle_position').textContent = bind_data['throttle_position']
    document.getElementById('text_relative_accelerator_pedal_position').textContent = bind_data['relative_accelerator_pedal_position']
    document.getElementById('text_engine_fuel_rate').textContent = bind_data['engine_fuel_rate']
    document.getElementById('text_absolute_barometric_pressure').textContent = bind_data['absolute_barometric_pressure']
    document.getElementById('text_fuel_pressure').textContent = bind_data['fuel_pressure']
    document.getElementById('text_service_distance').textContent = bind_data['service_distance']
    document.getElementById('text_anti_lock_barking_active').textContent = bind_data['anti_lock_barking_active']
    document.getElementById('text_steering_wheel_engine').textContent = bind_data['steering_wheel_engine']
    document.getElementById('text_position_of_doors').textContent = bind_data['position_of_doors']
    document.getElementById('text_right_left_turn_signal_light').textContent = bind_data['right_left_turn_signal_light']
    document.getElementById('text_alternate_beam_head_light').textContent = bind_data['alternate_beam_head_light']
    document.getElementById('text_high_beam_head_light').textContent = bind_data['high_beam_head_light']

    VehicleConfig.data.datasets[0].value = bind_data['vehicle_speed'];
    EngineConfig.data.datasets[0].value = bind_data['engine_speed'];
    barChartData.data.datasets[0].data = [bind_data['engine_coolant_temperature'], bind_data['intake_air_temperature'], bind_data['ambient_air_temperature'], bind_data['engine_oil_temperature']];
    
    window.VehicleGauge.update();
    window.engineGauge.update();
    window.barChart.update();
}

async function rending_dashboard() {
    var datasets = await get_dashboard_data()
    if (datasets) {
        hideSpinner()
        document.getElementById("dashboard_body").style.visibility = "visible";
        update_dashboard_components(datasets)
    }
    obd2_datasets = datasets
}

// Function to hide the Spinner 
function hideSpinner() {
    document.getElementById("spinner").style.visibility = "hidden";
}


function update_line_chart() {
    var x = document.getElementById("dropdown-menu").value;
    var json_datasets = []
    var timestamp = []

    for (var i = 0; i < obd2_datasets.length; i++) {
        json_datasets.push(obd2_datasets[i].data)
        timestamp.push(obd2_datasets[i].timestamp * 1000)
    }

    lineChartData.data.labels = timestamp

    lineChartData.data.datasets[0].data = json_datasets.map(function (e) {
        return e[x]
    })

    window.lineChart.update()
}

function update_drop_down_labels(data) {
    var labels = Object.keys(data)
    // Only need to update dropdown menu at first time
    for (var i = 0; i < labels.length; i++) {
        $('#dropdown-menu').append(' <option value="' + labels[i] + '">' + labels[i].replace(/_/g, ' ') + '</option>');
    }
}

rending_dashboard();

setInterval(async function () {
    var datasets = await get_dashboard_data();
    update_dashboard_components(datasets)
}, 60000); // Update every 60s
