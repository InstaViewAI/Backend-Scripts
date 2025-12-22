package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func getUpdateMap() map[string]map[string]any {
	jsonStr := `{
  "0x0003": {
    "0x0003:0x00": {
      "description": "Indicates the version of the device's cluster group. Used to manage compatibility across devices.\n指示设备群集组的版本。用于管理设备间的兼容性。",
      "name": "Identity\n身份"
    },
    "0x0003:0x01": {
      "description": "Specifies the currently installed firmware version on the device.\n指定设备上当前安装的固件版本。",
      "name": "Firmware Version\n固件版本"
    },
    "description": "Contains attributes related to the unique identity of the device.\n包含与设备唯一身份相关的属性。",
    "name": "Identity\n身份"
  },
  "0x0030": {
    "0x0030:0x01": {
      "description": "Values needed by firmware to know hardware properties\n固件识别硬件属性所需的值",
      "name": "Specifications\n规格"
    },
    "description": "Values that define the hardware configuration\n定义硬件配置的值",
    "name": "Hardware Configurations\n硬件配置"
  },
  "0x0032": {
    "0x0032:0x00": {
      "description": "Sets the verbosity of the device logs. Higher levels provide more detail for troubleshooting, while lower levels focus on essential system events.\n设置设备日志的详细级别。较高级别提供更多调试细节，较低级别关注关键系统事件。",
      "name": "Log level\n日志级别"
    },
    "description": "Controls how much system activity the camera logs, from detailed debugging to critical errors. Useful for troubleshooting and monitoring device behavior.\n控制摄像机记录系统活动的详细程度，从详细调试到关键错误。对故障排除和监控设备行为很有用。",
    "name": "Device Log Level\n设备日志级别"
  },
  "0xFC00": {
    "0xFC00:0x00": {
      "description": "Enables or disables privacy mode. When enabled, the camera physically moves to hide its lens and stop video recording.\n启用或禁用隐私模式。启用时，摄像机会物理移动以遮挡镜头并停止录像。",
      "name": "Privacy Mode\n隐私模式"
    },
    "description": "Controls the camera’s physical privacy mechanism by moving the lens or body to block video capture, ensuring complete visual privacy.\n通过移动镜头或机身来控制摄像机的物理隐私机制，阻止视频拍摄，确保完全的视觉隐私。",
    "name": "Privacy Mode\n隐私模式"
  },
  "0xFC01": {
    "0xFC01:0x00": {
      "description": "Sets the motion sensitivity level—Low, Medium, or High—based on how much movement is required to trigger detection.\n设置运动灵敏度等级—低、中、高—基于触发检测所需的动作量。",
      "name": "Motion Sensing Sensitivity Level\n运动感应灵敏度等级"
    },
    "description": "Controls how much motion is needed to trigger a detection event. Higher sensitivity means smaller movements will be detected.\n控制触发检测事件所需的运动量。灵敏度越高，越小的动作都能被检测到。",
    "name": "Motion Sensing Sensitivity Level\n运动感应灵敏度等级"
  },
  "0xFC02": {
    "0xFC02:0x00": {
      "description": "Image Flip Mode\n图像翻转模式",
      "name": "Image Flip Mode\n图像翻转模式"
    },
    "description": "Controls flipping of the camera’s video stream to match mounting orientation (e.g., ceiling or wall).\n控制摄像机视频流的翻转以匹配安装方向（如天花板或墙壁）。",
    "name": "Image Flip\n图像翻转"
  },
  "0xFC03": {
    "0xFC03:0x00": {
      "description": "Status Light\n状态灯",
      "name": "Status Light\n状态灯"
    },
    "description": "Turns the camera’s status light on or off. Useful for indicating when the camera is connecting, active or offline.\n开启或关闭摄像机状态灯。用于指示摄像机连接中、活动或离线状态。",
    "name": "Status Light\n状态灯"
  },
  "0xFC04": {
    "0xFC04:0x00": {
      "description": "Specifies the time zone using a region-based identifier (e.g., 'America/New_York').\n使用基于地区的标识符指定时区（例如“America/New_York”）。",
      "name": "Offset\n偏移量"
    },
    "0xFC04:0x01": {
      "description": "Specifies the time zone offset from UTC (e.g., 'UTC+0').\n指定与UTC的时区偏移（例如“UTC+0”）。",
      "name": "Identifier\n标识符"
    },
    "description": "Sets the device's time zone to ensure accurate timestamps for events and recordings.\n设置设备时区，确保事件和录制的时间戳准确。",
    "name": "TimeZone\n时区"
  },
  "0xFC06": {
    "0xFC06:0x00": {
      "description": "Light\n灯光",
      "name": "Light\n灯光"
    },
    "0xFC06:0x01": {
      "description": "Light Mode\n灯光模式",
      "name": "Light Mode\n灯光模式"
    },
    "description": "Controls the camera’s main spot light, which can operate in different modes such as infrared (IR), color (CLR), or intelligent (Intel) mode based on lighting conditions. 'Intel' automatically switches between infrared and color based on lighting conditions, 'IR' uses infrared for night vision, 'CLR' enables constant lighting for low-light clarity, and 'Off' disables the light entirely.\n控制摄像机的主聚光灯，能根据光照条件在红外（IR）、彩色（CLR）或智能（Intel）模式间切换。'Intel'根据光照自动切换红外和彩色，'IR'使用红外夜视，'CLR'保持常亮以提高清晰度，'Off'完全关闭灯光。",
    "name": "Camera Light\n摄像机灯光"
  },
  "0xFC07": {
    "0xFC07:0x00": {
      "description": "Sets the PIR detection sensitivity—Low, Medium, or High—based on the level of heat or motion needed to trigger detection.\n设置PIR检测灵敏度—低、中、高—基于触发检测所需的热量或运动水平。",
      "name": "PIR Sensitivity Level\nPIR灵敏度等级"
    },
    "description": "Adjusts how sensitive the camera’s PIR (Passive Infrared) sensor is to heat and motion. Higher levels detect smaller or farther movements.\n调节摄像机PIR（被动红外）传感器对热量和运动的灵敏度。级别越高能检测到更小或更远的动作。",
    "name": "PIR Sensor Sensitivity Level\nPIR传感器灵敏度等级"
  },
  "0xFC09": {
    "0xFC09:0x00": {
      "description": "Night Vision Mode\n夜视模式",
      "name": "Night Vision Mode\n夜视模式"
    },
    "description": "Controls the night vision mode and use of infrared light during low-light conditions.'Auto' enables infrared light automatically in low light, 'On' keeps infrared night vision always active, and 'Off' disables it entirely.\n控制夜视模式及低光条件下红外灯的使用。'Auto'在低光自动开启红外灯，'On'始终开启红外夜视，'Off'完全关闭。",
    "name": "Night Vision\n夜视"
  },
  "0xFC0A": {
    "0xFC0A:0x00": {
      "description": "IV Logo\nIV徽标",
      "name": "IV Logo\nIV徽标"
    },
    "0xFC0A:0x01": {
      "description": "Time\n时间",
      "name": "Time\n时间"
    },
    "description": "Controls what information is displayed on the video stream, such as time and logo overlays.\n控制视频流上显示的信息，如时间和徽标叠加。",
    "name": "OSD Settings\nOSD设置"
  },
  "0xFC0B": {
    "0xFC0B:0x00": {
      "description": "Microphone Enabled\n麦克风启用",
      "name": "Microphone Enabled\n麦克风启用"
    },
    "0xFC0B:0x01": {
      "description": "Speaker Enabled\n扬声器启用",
      "name": "Speaker Enabled\n扬声器启用"
    },
    "0xFC0B:0x02": {
      "description": "Volume Level\n音量级别",
      "name": "Volume Level\n音量级别"
    },
    "description": "Controls the camera’s audio functionality, including microphone, speaker, and volume level.\n控制摄像机的音频功能，包括麦克风、扬声器和音量级别。",
    "name": "Audio Settings\n音频设置"
  },
  "0xFC0C": {
    "0xFC0C:0x00": {
      "description": "Enabled\n启用",
      "name": "Enabled\n启用"
    },
    "0xFC0C:0x01": {
      "description": "Sets the time gap, in seconds, between two motion events. A lower value means more frequent events, while a higher value reduces event frequency.\n设置两个运动事件之间的时间间隔（秒）。值越小事件越频繁，值越大事件越少。",
      "name": "Cool Down Period\n冷却时间"
    },
    "description": "Controls whether motion detection is active and how frequently alerts are triggered by adjusting the cool down period between motion events.\n控制是否启用运动检测以及通过调整运动事件之间的冷却时间来决定警报触发的频率。",
    "name": "Motion Detection\n运动检测"
  },
  "0xFC0D": {
    "0xFC0D:0x00": {
      "description": "Turns tracking on or off. When enabled, the camera follows detected movement, either all motion or only human movement, based on camera capability.\n开启或关闭跟踪。启用时，摄像机会跟踪检测到的运动，可能是所有运动或仅人体运动，取决于摄像机功能。",
      "name": "Motion Tracking\n运动跟踪"
    },
    "description": "Enables tracking of movement within the camera’s view. Depending on the camera's capability, this may include general motion or only human motion tracking.\n启用摄像机视野内的运动跟踪。根据摄像机功能，这可能包括一般运动或仅限人体运动跟踪。",
    "name": "Human/Motion Tracking\n人体/运动跟踪"
  },
  "0xFC0F": {
    "0xFC0F:0x01": {
      "description": "Turns event scheduling on or off. When disabled, events are recorded at all times.\n开启或关闭事件调度。禁用时，全天记录事件。",
      "name": "Enabled\n启用"
    },
    "0xFC0F:0x02": {
      "description": "Specifies the start of the scheduled event recording window, in minutes from midnight (0–1439).\n指定事件调度记录窗口的开始时间，单位为从午夜起的分钟（0-1439）。",
      "name": "StartTime\n开始时间"
    },
    "0xFC0F:0x03": {
      "description": "Specifies the end of the scheduled event recording window, in minutes from midnight (0–1439). Must be equal to or greater than StartTime.\n指定事件调度记录窗口的结束时间，单位为从午夜起的分钟（0-1439）。必须大于或等于开始时间。",
      "name": "EndTime\n结束时间"
    },
    "description": "Defines a daily time window during which motion events will be recorded. Useful for limiting event detection to specific hours of the day.\n定义每天的时间窗口，在此期间记录运动事件。用于限制事件检测在特定时间段内进行。",
    "name": "EventScheduling\n事件调度"
  },
  "0xFC11": {
    "0xFC11:0x00": {
      "description": "Enabled\n启用",
      "name": "Enabled\n启用"
    },
    "0xFC11:0x01": {
      "description": "Zones\n区域",
      "name": "Zones\n区域"
    },
    "description": "Defines areas in the camera’s view where motion events should be ignored. These zones help reduce unwanted alerts by blocking motion detection.\n定义摄像机视野中应忽略运动事件的区域。通过屏蔽运动检测，这些区域有助于减少不必要的警报。",
    "name": "Block Activity Zone\n阻止活动区域"
  },
  "0xFC12": {
    "0xFC12:0x00": {
      "description": "Pan\n水平旋转",
      "name": "Pan\n水平旋转"
    },
    "0xFC12:0x01": {
      "description": "Tilt\n垂直倾斜",
      "name": "Tilt\n垂直倾斜"
    },
    "description": "Controls Pan, Tilt, and Zoom (PTZ) settings for the camera.\n控制摄像机的云台（PTZ）设置，包括水平旋转和垂直倾斜。",
    "name": "PTZ\n 云台"
  },
  "0xFC13": {
    "0xFC13:0x00": {
      "description": "Specifies the maximum length of a motion event recording (10–300 seconds). If motion stops earlier, the recording ends automatically to save battery.\n指定运动事件录制的最长时间（10-300秒）。若运动提前停止，录制自动结束以节省电池。",
      "name": "Event Duration\n事件时长"
    },
    "description": "Defines the maximum duration of each event for motion-triggered recordings on battery-powered cameras, in seconds.\n定义电池供电摄像机中每个运动触发录制事件的最长时长（秒）。",
    "name": "Event Duration\n事件时长"
  },
  "0xFC14": {
    "0xFC14:0x00": {
      "description": "When enabled, the camera’s CPU stays active at all times for faster wake-up and processing. When disabled, the CPU sleeps between events to conserve battery.\n启用时，摄像机CPU始终活跃，加快唤醒和处理速度。禁用时，事件间CPU休眠以节省电池。",
      "name": "Always ON\n始终开启"
    },
    "description": "Controls whether the CPU of a battery-powered camera remains always on.\n控制电池供电摄像机的CPU是否始终开启。",
    "name": "Always ON\n始终开启"
  },
  "0xFC15": {
    "0xFC15:0x00": {
      "description": "Current arm session identifier.\n当前布防会话标识符。",
      "name": "Arm ID\n布防ID"
    },
    "0xFC15:0x01": {
      "description": "Exit Delay. Time to wait after receiving arming command before actually arming the camera.\n退出延迟。收到布防命令后等待的时间，然后实际布防摄像机。",
      "name": "Exit Delay\n退出延迟"
    },
    "0xFC15:0x02": {
      "description": "Security State. Current arm state of the camera.\n安全状态。摄像机当前的布防状态。",
      "name": "Security State\n安全状态"
    },
    "description": "Manages the camera’s pro-monitoring status, including current arming state, session ID, and exit delay. Helps control and track when and how the camera is armed for security.\n管理摄像机的专业监控状态，包括当前布防状态、会话ID和退出延迟。帮助控制和追踪摄像机何时以及如何布防以确保安全。",
    "name": "Security State\n安全状态"
  },
  "0xFC16": {
    "0xFC16:0x00": {
      "description": "Guide Lamp Mode\n指示灯模式",
      "name": "Guide Lamp Mode\n指示灯模式"
    },
    "description": "Controls the guide lamp on the camera, which helps with visibility at night to locate the device. 'Intel' turns on the light when motion is detected, 'CLR' keeps the light on for constant visibility, and 'Off' disables the light.\n控制摄像机的指示灯，帮助夜间定位设备。'Intel'在检测到运动时开启，'CLR'保持常亮，'Off'关闭灯光。",
    "name": "Guide Lamp\n指示灯"
  },
  "0xFC17": {
    "0xFC17:0x00": {
      "description": "Alarm Light\n报警灯",
      "name": "Alarm Light\n报警灯"
    },
    "description": "Controls the alarm light on the camera, which flashes when motion is detected.\n控制摄像机的报警灯，在检测到运动时闪烁。",
    "name": "Alarm Light\n报警灯"
  },
  "0xFC18": {
    "0xFC18:0x00": {
      "description": "Flood Light Mode\n泛光灯模式",
      "name": "Flood Light Mode\n泛光灯模式"
    },
    "description": "Controls the flood light on the camera, used to illuminate the area for visibility, deterrence, or color video recording at night. 'Intel' turns it on automatically based on lighting conditions \u0026 motion, 'CLR' keeps it always on, and 'Off' disables it.\n控制摄像机的泛光灯，用于照亮区域以提高可见性、威慑或夜间彩色视频录制。'Intel'根据光照和运动自动开启，'CLR'常亮，'Off'关闭。",
    "name": "Flood Light\n泛光灯"
  },
  "0xFC19": {
    "0xFC19:0x00": {
      "description": "Enables or disables HDR mode. When enabled, the camera captures more detail in scenes with both bright and dark areas.\n启用或禁用HDR模式。启用时，摄像机能捕捉亮暗兼备场景中的更多细节。",
      "name": "HDR Enable\n启用HDR"
    },
    "description": "Controls whether High Dynamic Range (HDR) is enabled on the camera to improve image quality in challenging lighting conditions.\n控制摄像机是否启用高动态范围（HDR），以在复杂光照条件下提升图像质量。",
    "name": "HDR Enable\n启用HDR"
  },
  "0xFC1A": {
    "0xFC1A:0x00": {
      "description": "Defines the camera's bitrate mode. 'Auto' adjusts quality based on network conditions, while 'HD' prioritizes high video quality at a fixed rate.\n定义摄像机的码率模式。“Auto”根据网络状况调整质量，“HD”以固定码率优先保证高清视频质量。",
      "name": "BitRate\n 码率"
    },
    "description": "Controls the camera's video bitrate setting, affecting video quality and bandwidth usage.\n控制摄像机视频码率设置，影响视频质量和带宽使用。",
    "name": "BitRate\n 码率"
  },
  "0xFC1B": {
    "0xFC1B:0x00": {
      "description": "Enables or disables Continuous Video Recording to the SD card. When enabled, the camera records non-stop footage.\n启用或禁用SD卡上的连续视频录制。启用时，摄像机连续录像。",
      "name": "CVR"
    },
    "description": "Controls SD card-related settings for the camera, such as continuous video recording (CVR).\n控制摄像机的SD卡相关设置，例如连续视频录制（CVR）。",
    "name": "SD Card Settings\n SD卡设置"
  },
  "0xFC1C": {
    "0xFC1C:0x00": {
      "description": "Enables on-device AI to detect and classify people in motion events.\n启用设备端AI检测并分类运动事件中的人物。",
      "name": "Person\n人"
    },
    "0xFC1C:0x01": {
      "description": "Enables on-device AI to detect and classify animals in motion events.\n启用设备端AI检测并分类运动事件中的动物。",
      "name": "Animal\n动物"
    },
    "0xFC1C:0x02": {
      "description": "Enables on-device AI to detect and classify vehicles in motion events.\n启用设备端AI检测并分类运动事件中的车辆。",
      "name": "Vehicle\n车辆"
    },
    "description": "Defines which on-device AI models are enabled to detect specific types of motion events, such as people, animals, or vehicles.\n定义启用哪些设备端AI模型以检测特定类型的运动事件，如人、动物或车辆。",
    "name": "Edge AI\n 边缘AI"
  },
  "0xFC1D": {
    "0xFC1D:0x00": {
      "description": "Specifies the type of siren sound to play when triggered manually.\n指定手动触发时播放的警报声音类型。",
      "name": "Manual Siren\n 手动警报"
    },
    "description": "Triggers a siren sound manually from the camera to deter intruders or draw attention.\n手动触发摄像机发出警报声以驱赶入侵者或吸引注意。",
    "name": "Manual Siren\n 手动警报"
  },
  "0xFC1E": {
    "0xFC1E:0x00": {
      "description": "Indicates the Received Signal Strength Indicator (RSSI) in dBm.\n指示接收信号强度指示器（RSSI），单位为dBm。",
      "name": "RSSI"
    },
    "description": "Reports the camera’s current Wi-Fi signal strength (RSSI) in dBm, helping users assess network quality for optimal placement.\n报告摄像机当前的Wi-Fi信号强度（RSSI，单位dBm），帮助用户评估网络质量以优化放置位置。",
    "name": "RSSI"
  }
}
`

	// Declare a variable to hold the unmarshaled data
	var data map[string]map[string]any

	// Unmarshal JSON string to the map
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Use the map
	fmt.Printf("%+v\n", data)

	return data

}
