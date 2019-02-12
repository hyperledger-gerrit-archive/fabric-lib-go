// Copyright IBM Corp All Rights Reserved
//
// SPDX-License-Identifier: Apache-2.0
//
@Library("fabric-ci-lib") _
timeout(40) {
timestamps {
node ('snd-ubuntu1604-s390x-2c-16g-22') { // trigger build on x86_64 node

 stage("Clean Environment") {
   wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'xterm']) {
      commonSetup.cleanup() // Cleanup the leftover build artifats
      commonSetup.output() // Show the Jenkins environment details on console output
   }
 }
    try {
     def ROOTDIR = pwd() // workspace dir (/w/workspace/<job_name>)
     def PROJECT_DIR = "${BASE_DIR}"
// delete working directory
     deleteDir()
     commonSetup.cloneRepo 'fabric-lib-go'
        props = commonSetup.loadProperties()
        /* def props = ""
        dir("${ROOTDIR}/$PROJECT_DIR") {
          props = readProperties file:  'ci.properties'
          println props["GO_VER"]
        }
*/
       

        env.GOROOT = "/opt/go/go" + props["GO_VER"] + ".linux." + props["ARCH"]
        env.GOPATH = "$WORKSPACE/gopath"
        env.PATH = "$GOROOT/bin:$GOPATH/bin:$PATH"
        println "${PATH}"

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
