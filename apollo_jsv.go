/*
Copyright (c) 2013, 2014, Daniel Gruber (dgruber@univa.com), Univa

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"fmt"
	"github.com/dgruber/jsv"
	"regexp"
	"strconv"
	//	"strings"
)

func jsv_on_start_function() {
	//jsv.JSV_send_env()
}

func job_verification_function() {
	// Can be used for displaying submission parameters and
	// submission environment variables.
	jsv.JSV_log_info("--------------- Initial Params -----------------")
	jsv.JSV_show_params()
	jsv.JSV_show_envs()
	jsv.JSV_log_info("------------------------------------------------")

	// Automatic Core Binding
	// -----------------------------------------------------------------------------------------------------
	// Check if any binding has been already defined by the user.
	// If nothing has been set, and the user is in one of the defined parallel environments
	// then set an automatic linear core-binding set to the (min) number of slots requested
	slots := "1"
	if binding_type, binding_exists := jsv.JSV_get_param("binding_type"); binding_exists {
		jsv.JSV_log_info(fmt.Sprintf("Binding type: %s already set. No automatic core binding", binding_type))
	} else {
		pe_name, pe_exists := jsv.JSV_get_param("pe_name")
		if pe_exists {
			// Compile regex to match all PEs that we want to do automatic
			// core-binding on
			validPE := regexp.MustCompile(`(openmp|smp)`)
			if validPE.MatchString(pe_name) {
				pe_min, _ := jsv.JSV_get_param("pe_min")
				slots = pe_min
				jsv.JSV_log_info(fmt.Sprintf("Parallel Env: %s qualifies for automatic core binding", pe_name))
			} else {
				jsv.JSV_log_info(fmt.Sprintf("Parallel Env: %s does not qualify for automatic core binding", pe_name))
			}
		}
		// setting -binding linear:1 to each job (so that each
		// job can only use one core on the compute node)
		jsv.JSV_log_info(fmt.Sprintf("Setting automatic core binding to: %s cores", slots))
		jsv.JSV_set_param("binding_strategy", "linear_automatic")
		jsv.JSV_set_param("binding_type", "set")
		jsv.JSV_set_param("binding_amount", slots)
		jsv.JSV_set_param("binding_exp_n", "0")
	}

	// Automatic Memory Limitation
	// -----------------------------------------------------------------------------------------------------
	// Check if any memory limit has been set by the user.
	// If nothing has been set, define a memory limit as a multiple of the number of slots requested.
	// Use hard-coded value of 2GiB*slots.
	const memory_per_slot int = 2048
	if memory, exists := jsv.JSV_sub_get_param("l_hard", "m_mem_free"); exists {
		jsv.JSV_log_info(fmt.Sprintf("Memory limit of: %s already set. No automatic memory limiting", memory))
	} else {
		// convert slots to an integer, then calculate memory limit
		if num_slots, err := strconv.Atoi(slots); err == nil {
			memory_limit := num_slots * memory_per_slot
			jsv.JSV_log_info(fmt.Sprintf("Setting automatic memory limit of: %d MiB", memory_limit))
			jsv.JSV_sub_add_param("l_hard", "m_mem_free", strconv.Itoa(memory_limit))
		} else {
			jsv.JSV_log_error(fmt.Sprintf("Unable to convert: %s to an integer. Unable to set memory limit", slots))
		}
	}

	// Time Limitation
	jsv.JSV_sub_add_param("l_hard", "h_rt", strconv.Itoa(3600))

	// Can be used with server side JSV script to log
	// in qmaster messages file. For client side JSV
	// scripts to print out some messages when doing
	// qsub.
	//jsv.JSV_log_info("info message")
	//jsv.JSV_log_warning("warning message")
	//jsv.JSV_log_error("error message")

	// accepting the job but indicating that we did
	// some changes
	jsv.JSV_log_info("--------------- Final Params -----------------")
	jsv.JSV_show_params()
	jsv.JSV_log_info("------------------------------------------------")
	jsv.JSV_correct("Job was modified")
	return
}

/* example JSV 'script' */
func main() {
	jsv.Run(true, job_verification_function, jsv_on_start_function)
}
