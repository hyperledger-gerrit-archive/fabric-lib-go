// Copyright IBM Corp All Rights Reserved
//
// SPDX-License-Identifier: Apache-2.0
//
@Library("fabric-ci-lib") _
timeout(40) {
node ('hyp-x') { // trigger build on x86_64 node
 commonSetup.cleanup() // Cleanup the leftover build artifats
 commonSetup.output() // Show the Jenkins environment details on console output
 timestamps {
    try {
     def ROOTDIR = pwd() // workspace dir (/w/workspace/<job_name>)
  //   env.GOPATH = "$WORKSPACE/gopath"
// delete working directory
     deleteDir()
     commonSetup.cloneRepo 'fabric-lib-go'
// Run Unit-Tests
      stage("Unit Tests") {
         wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'xterm']) {
           try {
                 dir("${ROOTDIR}/$PROJECT_DIR/fabric-lib-go") {
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
