package main

import (
	"fmt"
	"github.com/dgruber/jsv"
	"regexp"
	//	"strconv"
	//	"strings"
)

func jsv_on_start_function() {
	//jsv.JSV_send_env()
}

func job_verification_function() {

	job_id, _ := jsv.JSV_get_param("JOB_ID")

	if context, exists := jsv.JSV_get_param("ac"); exists == false || context != "jsv" {
		jsv.JSV_accept("Job was accepted")
		return
	}

	// For certain queues we don't want the JSV to take effect
	queues_to_ignore := regexp.MustCompile(`admin.q.*`)
	if queue, exists := jsv.JSV_get_param("q_hard"); exists && queues_to_ignore.MatchString(queue) == true {
		jsv.JSV_accept("Job was accepted")
		return
	}

	// Can be used for displaying submission parameters and
	// submission environment variables.
	//jsv.JSV_log_info("--------------- Initial Params -----------------")
	//jsv.JSV_show_params()
	//jsv.JSV_log_info("------------------------------------------------")

	// Automatic Core Binding
	// -----------------------------------------------------------------------------------------------------
	// Check if any binding has been already defined by the user.
	// If nothing has been set, and the user is in one of the defined parallel environments
	// then set an automatic linear core-binding set to the (min) number of slots requested
	if binding_type, binding_exists := jsv.JSV_get_param("binding_type"); binding_exists {
		jsv.JSV_log_info(fmt.Sprintf("Binding type: %s already set. No automatic core binding", binding_type))
	} else if pe_name, pe_exists := jsv.JSV_get_param("pe_name"); pe_exists {
		// Compile regex to match all PEs that we want to do automatic
		// core-binding on
		validPE := regexp.MustCompile(`^(openmp|smp)$`)
		if validPE.MatchString(pe_name) {
			pe_min, _ := jsv.JSV_get_param("pe_min")
			slots := pe_min
			jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Parallel Env: %s qualifies for automatic core binding", job_id, pe_name))
			// setting -binding linear:slots to each job binding the job to the same number of cores as slots
			jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Setting automatic core binding to: %s cores", job_id, slots))
			jsv.JSV_set_param("binding_strategy", "linear_automatic")
			jsv.JSV_set_param("binding_type", "set")
			jsv.JSV_set_param("binding_amount", slots)
			jsv.JSV_set_param("binding_exp_n", "0")
		} else {
			jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Parallel Env: %s does not qualify for automatic core binding", job_id, pe_name))
		}
	} else {
		// No parallel environment set. So bind to 1 core only.
		// setting -binding linear:1 to each job (so that each
		// job can only use one core on the compute node)
		jsv.JSV_log_info(fmt.Sprintf("Job ID: %s. Setting automatic core binding to: %s cores", job_id, "1"))
		jsv.JSV_set_param("binding_strategy", "linear_automatic")
		jsv.JSV_set_param("binding_type", "set")
		jsv.JSV_set_param("binding_amount", "1")
		jsv.JSV_set_param("binding_exp_n", "0")
	}

	// Automatic Memory Limitation
	// -----------------------------------------------------------------------------------------------------
	// Check if any memory limit has been set by the user.
	// If nothing has been set, define a memory limit as a multiple of the number of slots requested.
	// Use hard-coded value of 2GiB*slots.
	//	const memory_per_slot int = 2048
	//	if memory, exists := jsv.JSV_sub_get_param("l_hard", "m_mem_free"); exists {
	//		jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Memory limit of: %s already set. No automatic memory limiting", job_id, memory))
	//	} else {
	//		// convert slots to an integer, then calculate memory limit
	//		if num_slots, err := strconv.Atoi(slots); err == nil {
	//			memory_limit := num_slots * memory_per_slot
	//			jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Setting automatic memory limit of: %d MiB", job_id, memory_limit))
	//			jsv.JSV_sub_add_param("l_hard", "m_mem_free", strconv.Itoa(memory_limit)+"M")
	//		} else {
	//			jsv.JSV_log_error(fmt.Sprintf("Unable to convert: %s to an integer. Unable to set memory limit", slots))
	//		}
	//	}

	// Time Limitation
	//	if h_rt, exists := jsv.JSV_sub_get_param("l_hard", "h_rt"); exists {
	//		jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Time limit of: %s already set. No automatic time limiting", job_id, h_rt))
	//	} else if h_rt, exists := jsv.JSV_sub_get_param("l_hard", "{~}h_rt"); exists {
	//		jsv.JSV_log_info(fmt.Sprintf("Job ID: %s Time limit of: %s already set. No automatic time limiting", job_id, h_rt))
	//	} else {
	//		jsv.JSV_reject("No job time limit set. Please either set a time limit or choose a job class from: default, default.medium, default.long")
	//		return
	//	}

	// accepting the job but indicating that we did
	// some changes
	//jsv.JSV_log_info("--------------- Final Params -----------------")
	//jsv.JSV_show_params()
	//jsv.JSV_log_info("------------------------------------------------")
	jsv.JSV_correct("Job was modified")
	return
}

/* example JSV 'script' */
func main() {
	jsv.Run(true, job_verification_function, jsv_on_start_function)
}
