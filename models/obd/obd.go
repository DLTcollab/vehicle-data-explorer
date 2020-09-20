package obd

type ODB2_data struct {
	Vin                                 string  `json:"vin"`
	Engine_load                         int     `json:"engine_load"`
	Engine_coolant_temperature          int     `json:"engine_coolant_temperature"`
	Fuel_pressure                       int     `json:"fuel_pressure"`
	Engine_speed                        float32 `json:"engine_speed"`
	Vehicle_speed                       int     `json:"vehicle_speed"`
	Intake_air_temperature              int     `json:"intake_air_temperature"`
	Mass_air_flow                       int     `json:"mass_air_flow"`
	Fuel_tank_level_input               int     `json:"fuel_tank_level_input"`
	Absolute_barometric_pressure        int     `json:"absolute_barometric_pressure"`
	Control_module_voltage              float32 `json:"control_module_voltage"`
	Throttle_position                   int     `json:"throttle_position"`
	Ambient_air_temperature             int     `json:"ambient_air_temperature"`
	Relative_accelerator_pedal_position int     `json:"relative_accelerator_pedal_position"`
	Engine_oil_temperature              int     `json:"engine_oil_temperature"`
	Engine_fuel_rate                    float32 `json:"engine_fuel_rate"`
	Service_distance                    int     `json:"service_distance"`
	Anti_lock_barking_active            int     `json:"anti_lock_barking_active"`
	Steering_wheel_angle                int     `json:"steering_wheel_angle"`
	Position_of_doors                   int     `json:"position_of_doors"`
	Right_left_turn_signal_light        int     `json:"right_left_turn_signal_light"`
	Alternate_beam_head_light           int     `json:"alternate_beam_head_light"`
	High_beam_head_light                int     `json:"high_beam_head_light"`
}

type OBD struct {
	PID uint8
}

// Service 09
// PID:0x02
func get_vechicle_identification_number(data string) string {
	return data
}

// Service 01
// PID:0x04
func get_calculated_engine_load(data uint32) int {
	var A = data >> 24
	return int(A) * 100 / 255
}

// PID:0x05
func get_engine_coolant_temperature(data uint32) int {
	var A = data >> 24
	return int(A) - 40
}

// PID:0x0A
func get_fuel_pressure(data uint32) int {
	var A = data >> 24
	return int(A) * 3
}

// PID:0x0C
func get_engine_speed(data uint32) float32 {
	var A = data >> 24
	var B = data << 8
	B = B >> 24

	return (float32(A)*256 + float32(B)) / 4
}

// PID:0x0D
func get_vehicle_speed(data uint32) int {
	var A = data >> 24
	return int(A)
}

// PID:0x0F
func get_intake_air_temperature(data uint32) int {
	var A = data >> 24
	return int(A) - 40
}

// PID:0x10
func get_mass_air_flow(data uint32) int {
	var A = data >> 24
	var B = data << 8
	B = B >> 24

	return (int(A)*256 + int(B)) / 100
}

// PID:0x2F
func get_fuel_tank_level_input(data uint32) int {
	var A = data >> 24
	return int(A) * 100 / 255
}

// PID:0x33
func get_absolute_barometric_pressure(data uint32) int {
	var A = data >> 24
	return int(A)
}

// PID:0x42
func get_control_module_voltage(data uint32) float32 {
	var A = data >> 24
	var B = data << 8
	B = B >> 24
	var v = (float32(A)*256 + float32(B)) / 1000.0
	return v
}

// PID:0x45
func get_throttle_position(data uint32) int {
	var A = data >> 24
	return int(A) * 100 / 255
}

// PID:0x46
func get_ambient_air_temperature(data uint32) int {
	var A = data >> 24
	return int(A) - 40
}

// PID:0x5A
func get_relative_accelerator_pedal_position(data uint32) int {
	var A = data >> 24
	return int(A) * 100 / 255
}

// PID:0x5C
func get_engine_oil_temperature(data uint32) int {
	var A = data >> 24
	return int(A) - 40
}

// PID:0x5E
func get_engine_fuel_rate(data uint32) float32 {
	var A = data >> 24
	var B = data << 8
	B = B >> 24
	return (float32(A)*256 + float32(B)) / 20
}

// E1~E8 is not inside obd2-pids standard
// PID:E1
func get_service_distance(data uint32) uint32 {
	return data
}

// PID:E2
func get_anti_lock_barking_active(data uint32) uint32 {
	return data
}

// PID:E3
func get_steering_wheel_angle(data uint32) uint32 {
	return data
}

// PID:E5
func get_position_of_doors(data uint32) uint32 {
	return data
}

// PID:E6
func get_right_left_turn_signal_light(data uint32) uint32 {
	return data
}

// PID:E7
func get_alternate_beam_head_light(data uint32) uint32 {
	return data
}

// PID:E8
func get_high_beam_head_light(data uint32) uint32 {
	return data
}
