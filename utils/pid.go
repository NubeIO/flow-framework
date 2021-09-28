package utils

import (
	"time"
)

type PIDMode int

const (
	MANUAL = iota // 0
	AUTO          // 1
)

type PIDDirection int

const (
	DIRECT  = iota // 0
	REVERSE        // 1
)

type PIDController struct {
	pidInput      float64
	pidEnable     bool
	pidSP         float64
	pidKp         float64 // * (P)roportional Tuning Parameter
	pidKi         float64 // * (I)ntegral Tuning Parameter
	pidKd         float64 // * (D)erivative Tuning Parameter
	pidOutput     float64
	pidOutMax     float64
	pidOutMin     float64
	pidMode       PIDMode
	pidDirection  PIDDirection
	pidBias       float64
	pidSampleTime int64
	pidLastInput  float64
	pidITerm      float64
	pidLastTime   int64
	pidDisplayKp  float64
	pidDisplayKi  float64
	pidDisplayKd  float64
}

func NewPIDController() PIDController {
	var pidInstance PIDController
	pidInstance.pidMode = MANUAL
	pidInstance.SetOutputLimits(0, 100) //default output limit
	pidInstance.pidSampleTime = 1000    //default Controller Sample Time is 1 seconds
	pidInstance.pidITerm = 0
	pidInstance.SetTunings(-1, -1, -1)
	pidInstance.SetControllerDirection(DIRECT)
	pidInstance.SetBias(0)
	pidInstance.pidLastTime = time.Now().UnixMilli() - pidInstance.pidSampleTime
	return pidInstance
}

/*Compute **********************************************************************
 *     This, as they say, is where the magic happens.  this function should be called
 *   every time "void loop()" executes.  the function will decide for itself whether a new
 *   pid Output needs to be computed.  returns true when the output is computed,
 *   false when nothing has been done.
 **********************************************************************************/
func (pid PIDController) Compute() bool {
	if pid.pidMode == MANUAL {
		return false
	}
	now := time.Now().UnixMilli()
	timeChange := now - pid.pidLastTime
	if timeChange >= pid.pidSampleTime {
		/*Compute all the working error variables*/
		input := pid.pidInput
		errorValue := input - pid.pidSP //input > setpoint = positive error
		errorValue = float64(pid.pidDirection) * errorValue

		pid.pidITerm += pid.pidKi * errorValue
		if pid.pidITerm > pid.pidOutMax-pid.pidBias {
			pid.pidITerm = pid.pidOutMax - pid.pidBias
		} else if pid.pidITerm < pid.pidOutMin-pid.pidBias {
			pid.pidITerm = pid.pidOutMin - pid.pidBias
		}

		dInput := input - pid.pidLastInput

		// Compute PID Output
		var output float64
		//output = ((this.kp * error) + this.ITerm - (this.kd * dInput)) * this.setDirection;
		output = pid.pidKp*errorValue + pid.pidITerm - pid.pidKd*dInput + pid.pidBias
		//output = ((this.kp * error) + this.ITerm - (this.kd * dInput));

		if output > pid.pidOutMax {
			output = pid.pidOutMax
		} else if output < pid.pidOutMin {
			output = pid.pidOutMin
		}
		pid.pidOutput = output

		// Remember some variables for next time
		pid.pidLastInput = input
		pid.pidLastTime = now
		return true
	} else {
		return false
	}
}

/*SetOutput ***************************************************************
 * Set output if in manual mode.
 *************************************************************************/
func (pid PIDController) SetOutput(outputValue float64) {
	if outputValue > pid.pidOutMax {
		outputValue = pid.pidOutMax // POSSIBLY INCORRECT
	} else if outputValue < pid.pidOutMin {
		outputValue = pid.pidOutMin
	}
	pid.pidOutput = outputValue
}

/*SetControllerDirection *************************************************
 * The PID will either be connected to a DIRECT acting process (+Input leads
 * to +Output) or a REVERSE acting process(-Input leads to +Output.)
 ******************************************************************************/
func (pid PIDController) SetControllerDirection(newDirection PIDDirection) {
	pid.pidDirection = newDirection
}

/*Initialize ****************************************************************
 *	does all the things that need to happen to ensure a bumpless transfer
 *  from manual to automatic mode.  NOT CURRENTLY USED.
 ******************************************************************************/
func (pid PIDController) Initialize() {
	//ITerm = PID_OUT;
	pid.pidITerm = 0
	pid.pidOutput = pid.pidBias
	pid.pidLastInput = pid.pidInput
	/*
		  if (pid.pidITerm > pid.pidOutMax) {
			pid.pidITerm = this.outMax;
		  } else if (pid.pidITerm < pid.pidOutMin) {
			pid.pidITerm = pid.pidOutMin;
		  }
	*/
}

/*SetMode ****************************************************************
 * Allows the controller Mode to be set to manual (0) or Automatic (non-zero)
 * when the transition from manual to auto occurs, the controller is
 * automatically initialized
 ******************************************************************************/
func (pid PIDController) SetMode(newMode PIDMode) {
	pid.pidMode = newMode
}

/*SetBias ****************************************************
 *     This function sets the Bias of the PID function
 **************************************************************************/
func (pid PIDController) SetBias(newBias float64) {
	if newBias > pid.pidOutMax {
		newBias = pid.pidOutMax // POSSIBLY INCORRECT
	} else if newBias < pid.pidOutMin {
		newBias = pid.pidOutMin
	}
	pid.pidBias = newBias
}

/*SetOutputLimits ****************************************************
 * Clamps the output values between limits.
 **************************************************************************/
func (pid PIDController) SetOutputLimits(Min float64, Max float64) {
	if Min >= Max {
		return
	}
	pid.pidOutMin = Min
	pid.pidOutMax = Max
	if pid.pidMode == AUTO {
		if pid.pidOutput > pid.pidOutMax {
			pid.pidOutput = pid.pidOutMax
		} else if pid.pidOutput < pid.pidOutMin {
			pid.pidOutput = pid.pidOutMin
		}

		if pid.pidITerm > pid.pidOutMax-pid.pidBias {
			pid.pidITerm = pid.pidOutMax - pid.pidBias
		} else if pid.pidITerm < pid.pidOutMin-pid.pidBias {
			pid.pidITerm = pid.pidOutMin - pid.pidBias
		}
	}
}

/*SetSampleTime *********************************************************
 * sets the period, in Milliseconds, at which the calculation is performed
 ******************************************************************************/
func (pid PIDController) SetSampleTime(NewSampleTime int64) {
	if NewSampleTime > 0 {
		ratio := float64(NewSampleTime) / float64(pid.pidSampleTime)
		pid.pidKi *= ratio
		pid.pidKd /= ratio
		pid.pidSampleTime = NewSampleTime
	}
}

/*SetTunings  *************************************************************
 * This function allows the controller's dynamic performance to be adjusted.
 * it's called automatically from the constructor, but tunings can also
 * be adjusted on the fly during normal operation
 ******************************************************************************/
func (pid PIDController) SetTunings(Kp float64, Ki float64, Kd float64) {
	if Kp < 0 || Ki < 0 || Kd < 0 {
		return
	}
	pid.pidDisplayKp = Kp
	pid.pidDisplayKi = Ki
	pid.pidDisplayKd = Kd
	if Ki == 0 {
		pid.pidITerm = 0
	}
	SampleTimeInSec := float64(pid.pidSampleTime) / 1000
	pid.pidKp = Kp
	pid.pidKi = Ki * SampleTimeInSec
	pid.pidKd = Kd / SampleTimeInSec
}

/* Status Funcions*************************************************************
 * Just because you set the Kp=-1 doesn't mean it actually happened.  these
 * functions query the internal state of the PID.  they're here for display
 * purposes.
 ******************************************************************************/

func (pid PIDController) GetKp() float64             { return pid.pidDisplayKp }
func (pid PIDController) GetKi() float64             { return pid.pidDisplayKi }
func (pid PIDController) GetKd() float64             { return pid.pidDisplayKd }
func (pid PIDController) GetMode() PIDMode           { return pid.pidMode }
func (pid PIDController) GetDirection() PIDDirection { return pid.pidDirection }
func (pid PIDController) GetOutput() float64         { return pid.pidOutput }
func (pid PIDController) GetInput() float64          { return pid.pidInput }
func (pid PIDController) GetSetpoint() float64       { return pid.pidSP }
func (pid PIDController) GetBias() float64           { return pid.pidBias }

/*

class PID {
    public:

        //Constants used in some of the functions below
        typedef enum {
            MANUAL = 0,
            AUTO = 1,
        } PID_MODE;

        typedef enum {
            DIRECT = 0,
            REVERSE = 1,
        } PID_DIRECTION;


        //commonly used functions **************************************************************************
        PID();                                // * constructor.

        void SetMode(PID_MODE mode);               // * sets PID to either Manual (0) or Auto (non-0)
        void SetMode(bool mode);               // * sets PID to either Manual (0) or Auto (non-0)

        void Initialize();                    // does all the things that need to happen to ensure a bumpless
                                                // transfer from manual to automatic mode.

        float Compute();                       // * performs the PID calculation.  it should be
                                                //   called every time loop() cycles. ON/OFF and
                                                //   calculation frequency can be set using SetMode
                                                //   SetSampleTime respectively

        void SetOutput(float outputValue);   //   Sets the output value. To be used for manual control.

        void SetOutputLimits(float, float); // * clamps the output to a specific range. 0-255 by default, but
                                                //   it's likely the user will want to change this depending on
                                                //   the application



        //available but not commonly used functions ********************************************************
        void SetTunings(float, float,       // * While most users will set the tunings once in the
                        float);         	    //   constructor, this function gives the user the option
                                                //   of changing tunings during runtime for Adaptive control

        void SetBias(float);

        void SetControllerDirection(bool);
        void SetControllerDirection(PID_DIRECTION);	  // * Sets the Direction, or "Action" of the controller. DIRECT
                            //   means the output will increase when error is positive. REVERSE
                            //   means the opposite.  it's very unlikely that this will be needed
                            //   once it is set in the constructor.
        void SetSampleTime(int);              // * sets the frequency, in Milliseconds, with which
                                                //   the PID calculation is performed.  default is 100

        //Display functions ****************************************************************
        float GetKp();						  // These functions query the pid for interal values.
        float GetKi();						  //  they were created mainly for the pid front-end,
        float GetKd();						  // where it's important to know what is actually
        PID_MODE GetMode();						  //  inside the PID.
        PID_DIRECTION GetDirection();
        float GetOutput();
        float GetInput();
        float GetSetpoint();
        float GetBias();

        //IMPLEMENTATION functions ****************************************************************
        // All variables from Modbus or Physicl IO
        float implementPID(IOHandler* io_handler, PhysicalPointRef* PID_IN_POINT, Point PID_SP_POINT, Point
                PID_Px_POINT, Point PID_Ix_POINT, Point PID_Dx_POINT, Point PID_DIRECTION_POINT, Point PID_ENABLE_POINT);



    private:
        //IMPLEMENTATION VARIABLES
        float PID_IN;
        bool PID_ENABLE;
        float PID_SP;
        float PID_Px;
        float PID_Ix;
        float PID_Dx;
        float PID_OUT;

        float dispKp;				// * we'll hold on to the tuning parameters in user-entered
        float dispKi;				//   format for display purposes
        float dispKd;				//

        float kp;                  // * (P)roportional Tuning Parameter
        float ki;                  // * (I)ntegral Tuning Parameter
        float kd;                  // * (D)erivative Tuning Parameter

        float bias;

        unsigned long lastTime;
        float outputSum, lastInput, ITerm;

        unsigned long SampleTime;
        float outMin, outMax;
        PID_MODE currMode;
        PID_DIRECTION controllerDirection;

};
#endif
*/
