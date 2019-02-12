// Copyright IBM Corp All Rights Reserved
//
// SPDX-License-Identifier: Apache-2.0
//
@Library("fabric-ci-lib") _
timeout(40) {
node ('hyp-x') { // trigger build on x86_64 node
  timestamps {
  try {
  def ROOTDIR = pwd() // workspace dir (/w/workspace/<job_name>)
  def PROJECT_DIR = "${BASE_DIR}"
  stage("Setup Environment") {
   wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'xterm']) {
      fabBuildLibrary.cleanup() // Cleanup the leftover build artifats
      fabBuildLibrary.output() //  Show the Jenkins environment details on the console output
   }
  }
  stage("Clone Changes") {
      // Delete working directory
      deleteDir()
      commonSetup.cloneRepo 'fabric-lib-go'
  }
      props = commonSetup.loadProperties()
      // Set GOPATH
      env.GOROOT = "/opt/go/go" + props["GO_VER"] + ".linux." + props["ARCH"]
      env.GOPATH = "$WORKSPACE/gopath"
      env.PATH = "$GOROOT/bin:$GOPATH/bin"
      println "${PATH}"

// Run Checks
      stage("Unit Tests") {
         wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'xterm']) {
           try {
                 dir("${ROOTDIR}/$PROJECT_DIR") {
                 sh 'make checks'
                 }
               }
           catch (err) {
                 failure_stage = "checks"
                 currentBuild.result = 'FAILURE'
                 throw err
           }
         }
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
        commonSetup.rocketSend() // Send the merge build failure notifications to jenkins-robot rocketChat channel
        cleanWs()
      } // finally block
  } // timestamps block
} // node block block
} // timeout block