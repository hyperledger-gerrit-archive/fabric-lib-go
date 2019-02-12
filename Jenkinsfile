// Copyright IBM Corp All Rights Reserved
//
// SPDX-License-Identifier: Apache-2.0
//
@Library("fabric-ci-lib") _
timeout(40) {
node ('snd-ubuntu1604-s390x-2c-16g-22') { // trigger build on x86_64 node
 commonSetup.cleanup() // Cleanup the leftover build artifats
 commonSetup.output() // Show the Jenkins environment details on console output

 timestamps {
    try {
     def ROOTDIR = pwd() // workspace dir (/w/workspace/<job_name>)
     def PROJECT_DIR = "${BASE_DIR}"
// delete working directory
     deleteDir()
     commonSetup.cloneRepo 'fabric-lib-go'
        // commonSetup.loadProperties()
      // def props = readProperties file:'ci.properties';
      dir("${ROOTDIR}/$PROJECT_DIR") {
        def props = readFile 'ci.properties'
        echo "Immediate one ${props.GO_VER}"
        }
// Run Unit-Tests
      stage("Unit Tests") {
         wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'xterm']) {
           try {
                 dir("${ROOTDIR}/$PROJECT_DIR") {
                 sh 'make unit-tests'
                 }
               }
           catch (err) {
                 failure_stage = "Unit Tests"
                 currentBuild.result = 'FAILURE'
                 throw err
           }
         }
      }
    } finally { // Code for coverage report
    //    commonSetup.rocketSend() // Send the merge build failure notifications to jenkins-robot rocketChat channel
sh 'echo "ROCKETCHAT"'
      } // finally block
  } // timestamps block
} // node block block
} // timeout block
