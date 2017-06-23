/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package pipeline

/*
Controller is the interface that defines the methods for controlling a pipeline stage
*/
type Controller interface {
	Start() error
	Abort() error
	Quiesce() error
	Done()
}

/*
StageControl is a structure that contains an input and an output control
    channel to be used for controlling each stage of the pipeline
*/
type StageControl struct {
	Input  chan string
	Output chan string
}
