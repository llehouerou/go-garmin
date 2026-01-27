# Garmin Connect API Endpoints

This document lists all known API endpoints from the reference projects and web research:
- [python-garminconnect](https://github.com/cyberjunky/python-garminconnect)
- [garth](https://github.com/matin/garth)
- [garmin-workouts](https://github.com/mkuthan/garmin-workouts)
- [garmin-connect (JS)](https://github.com/Pythe1337N/garmin-connect)
- [dotnet.garmin.connect](https://github.com/sealbro/dotnet.garmin.connect)

## Implementation Status

- [x] Implemented
- [ ] Not implemented

---

## Sleep Service (`/sleep-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/sleep-service/sleep/dailySleepData?date={date}` | Daily sleep data |

---

## Wellness Service (`/wellness-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/wellness-service/wellness/dailyStress/{date}` | Daily stress and body battery data |
| [x] | GET | `/wellness-service/wellness/bodyBattery/events/{date}` | Body battery events |
| [x] | GET | `/wellness-service/wellness/dailyHeartRate/?date={date}` | Daily heart rate data |
| [ ] | GET | `/wellness-service/wellness/daily/spo2/{date}` | Daily SpO2 data |
| [ ] | GET | `/wellness-service/wellness/daily/respiration/{date}` | Daily respiration data |
| [ ] | GET | `/wellness-service/wellness/daily/im/{date}` | Daily intensity minutes |
| [ ] | GET | `/wellness-service/wellness/dailyEvents/{date}` | Daily events |
| [ ] | GET | `/wellness-service/wellness/dailySummaryChart/{displayName}?date={date}` | Daily summary chart |
| [ ] | GET | `/wellness-service/wellness/floorsChartData/daily/{date}` | Floor climbing data |
| [ ] | GET | `/wellness-service/wellness/epoch/request/{date}` | Request epoch data reload |
| [ ] | GET | `/wellness-service/wellness/bodyBattery/reports/daily?startDate={start}&endDate={end}` | Body battery reports |
| [ ] | GET | `/wellness-service/stats/daily/sleep/score/{start}/{end}` | Sleep score stats |

---

## Activity List Service (`/activitylist-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/activitylist-service/activities/search/activities?start={start}&limit={limit}` | Search activities |
| [ ] | GET | `/activitylist-service/activities/` | List activities |
| [ ] | GET | `/activitylist-service/activities/count` | Activity count |

---

## Activity Service (`/activity-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/activity-service/activity/{activityId}` | Get single activity |
| [ ] | GET | `/activity-service/activity/{activityId}/details` | Activity details |
| [x] | GET | `/activity-service/activity/{activityId}/splits` | Activity splits |
| [ ] | GET | `/activity-service/activity/{activityId}/typedSplits` | Activity typed splits |
| [ ] | GET | `/activity-service/activity/{activityId}/split_summaries` | Activity split summaries |
| [x] | GET | `/activity-service/activity/{activityId}/weather` | Activity weather |
| [ ] | GET | `/activity-service/activity/{activityId}/hrTimeInZones` | HR time in zones |
| [ ] | GET | `/activity-service/activity/{activityId}/powerTimeInZones` | Power time in zones |
| [ ] | GET | `/activity-service/activity/{activityId}/exerciseSets` | Exercise sets |
| [ ] | GET | `/activity-service/activity/{activityId}/gear` | Activity gear |
| [ ] | GET | `/activity-service/activity/activityTypes` | Activity types |
| [ ] | GET | `/activity-service/activity/forDate/{date}` | Activities for date |
| [ ] | DELETE | `/activity-service/activity/{activityId}` | Delete activity |

---

## Download Service (`/download-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/download-service/files/activity/{activityId}` | Download activity (original FIT) |
| [ ] | GET | `/download-service/export/tcx/activity/{activityId}` | Export activity as TCX |
| [ ] | GET | `/download-service/export/gpx/activity/{activityId}` | Export activity as GPX |
| [ ] | GET | `/download-service/export/kml/activity/{activityId}` | Export activity as KML |
| [ ] | GET | `/download-service/export/csv/activity/{activityId}` | Export activity as CSV |

---

## Upload Service (`/upload-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | POST | `/upload-service/upload` | Upload activity file (FIT, TCX, GPX) |

---

## Weight Service (`/weight-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/weight-service/weight/dayview/{date}` | Daily weight data |
| [ ] | GET | `/weight-service/weight/range/{start}/{end}?includeAll=true` | Weight range |
| [ ] | GET | `/weight-service/weight/dateRange?startDate={start}&endDate={end}` | Weight date range |
| [ ] | GET | `/weight-service/weight/daterangesnapshot` | Body composition snapshot |
| [ ] | DELETE | `/weight-service/weight/{date}/{weightPK}` | Delete weigh-in |

---

## User Summary Service (`/usersummary-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/usersummary-service/usersummary/daily/?calendarDate={date}` | Daily user summary |
| [ ] | GET | `/usersummary-service/usersummary/hydration/daily/{date}` | Daily hydration |
| [ ] | PUT | `/usersummary-service/usersummary/hydration/log` | Log/update hydration |
| [ ] | GET | `/usersummary-service/stats/steps/daily/{start}/{end}` | Daily steps stats |
| [ ] | GET | `/usersummary-service/stats/steps/weekly/{end}/{weeks}` | Weekly steps stats |
| [ ] | GET | `/usersummary-service/stats/stress/daily/{start}/{end}` | Daily stress stats |
| [ ] | GET | `/usersummary-service/stats/stress/weekly/{end}/{weeks}` | Weekly stress stats |
| [ ] | GET | `/usersummary-service/stats/hydration/daily/{start}/{end}` | Hydration stats |
| [ ] | GET | `/usersummary-service/stats/im/daily/{start}/{end}` | Daily intensity minutes |
| [ ] | GET | `/usersummary-service/stats/im/weekly/{start}/{end}` | Weekly intensity minutes |

---

## User Stats Service (`/userstats-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/userstats-service/wellness/daily/{date}` | Daily wellness stats (RHR) |

---

## HRV Service (`/hrv-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/hrv-service/hrv/{date}` | Daily HRV data |
| [x] | GET | `/hrv-service/hrv/daily/{start}/{end}` | HRV range |

---

## Metrics Service (`/metrics-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/metrics-service/metrics/trainingreadiness/{date}` | Training readiness |
| [ ] | GET | `/metrics-service/metrics/endurancescore?calendarDate={date}` | Endurance score |
| [ ] | GET | `/metrics-service/metrics/hillscore?calendarDate={date}` | Hill score |
| [ ] | GET | `/metrics-service/metrics/racepredictions?calendarDate={date}` | Race predictions |
| [ ] | GET | `/metrics-service/metrics/maxmet/daily/{start}/{end}` | Daily VO2 max/MET |
| [ ] | GET | `/metrics-service/metrics/maxmet/latest/{date}` | Latest VO2 max/MET |
| [ ] | GET | `/metrics-service/metrics/trainingstatus/aggregated` | Training status aggregated |
| [ ] | GET | `/metrics-service/metrics/trainingstatus/daily/{date}` | Daily training status |
| [ ] | GET | `/metrics-service/metrics/trainingloadbalance/latest/{date}` | Training load balance |
| [ ] | GET | `/metrics-service/metrics/heataltitudeacclimation/latest/{date}` | Heat/altitude acclimation |

---

## Device Service (`/device-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/device-service/deviceregistration/devices` | List devices |
| [ ] | GET | `/device-service/deviceservice/device-info/settings/{deviceId}` | Device settings |
| [ ] | GET | `/device-service/devicemessage/messages` | Device messages |

---

## Web Gateway (`/web-gateway/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/web-gateway/device-info/primary-training-device` | Primary training device |
| [ ] | GET | `/web-gateway/solar/{deviceId}?startDate={start}&endDate={end}` | Solar panel data |

---

## User Profile Service (`/userprofile-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/userprofile-service/socialProfile` | Social profile |
| [ ] | GET | `/userprofile-service/userprofile/user-settings` | User settings |
| [ ] | GET | `/userprofile-service/userprofile/settings` | Profile settings |
| [ ] | GET | `/userprofile-service/userprofile/profile` | User profile |

---

## Goal Service (`/goal-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/goal-service/goal/goals?status={status}` | Get goals |

---

## Gear Service (`/gear-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/gear-service/gear/filterGear?userProfilePk={pk}` | Get gear |
| [ ] | GET | `/gear-service/gear/stats/{gearUUID}` | Gear stats |
| [ ] | GET | `/gear-service/gear/activities/{gearUUID}?start={start}&limit={limit}` | Gear activities |

---

## Badge Service (`/badge-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/badge-service/badge/earned` | Earned badges |
| [ ] | GET | `/badge-service/badge/available` | Available badges |

---

## Badge Challenge Service (`/badgechallenge-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/badgechallenge-service/badgeChallenge/completed?start={start}&limit={limit}` | Completed badge challenges |
| [ ] | GET | `/badgechallenge-service/badgeChallenge/available?start={start}&limit={limit}` | Available badge challenges |
| [ ] | GET | `/badgechallenge-service/badgeChallenge/non-completed?start={start}&limit={limit}` | Non-completed badge challenges |
| [ ] | GET | `/badgechallenge-service/virtualChallenge/inProgress?start={start}&limit={limit}` | In-progress virtual challenges |

---

## Ad-Hoc Challenge Service (`/adhocchallenge-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/adhocchallenge-service/adHocChallenge/historical?start={start}&limit={limit}` | Historical ad-hoc challenges |

---

## Blood Pressure Service (`/bloodpressure-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/bloodpressure-service/bloodpressure/range/{start}/{end}` | Blood pressure range |
| [ ] | POST | `/bloodpressure-service/bloodpressure` | Log blood pressure |
| [ ] | DELETE | `/bloodpressure-service/bloodpressure/{version}/{date}` | Delete blood pressure |

---

## Personal Record Service (`/personalrecord-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/personalrecord-service/personalrecord/prs/{displayName}` | Personal records |

---

## Biometric Service (`/biometric-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/biometric-service/biometric/{displayName}?startDate={start}&endDate={end}` | Biometric data |
| [ ] | GET | `/biometric-service/stats/{displayName}?startDate={start}&endDate={end}` | Biometric stats |
| [ ] | GET | `/biometric-service/heartRateZones/{displayName}` | Heart rate zones |
| [ ] | GET | `/biometric-service/biometric/ftp` | Cycling FTP |
| [ ] | GET | `/biometric-service/biometric/lactatethreshold` | Lactate threshold |

---

## Fitness Age Service (`/fitnessage-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/fitnessage-service/fitnessage/{displayName}` | Fitness age |

---

## Fitness Stats Service (`/fitnessstats-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/fitnessstats-service/activity/all?startDate={start}&endDate={end}&...` | All activity stats |

---

## Workout Service (`/workout-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/workout-service/workouts?start={start}&limit={limit}` | List workouts |
| [ ] | GET | `/workout-service/workout/{workoutId}` | Get workout |
| [ ] | GET | `/workout-service/workout/FIT/{workoutId}` | Download workout FIT |
| [ ] | POST | `/workout-service/workout` | Create workout |
| [ ] | PUT | `/workout-service/workout/{workoutId}` | Update workout |
| [ ] | DELETE | `/workout-service/workout/{workoutId}` | Delete workout |
| [ ] | POST | `/workout-service/schedule/{workoutId}` | Schedule workout |
| [ ] | GET | `/workout-service/schedule/{scheduleId}` | Get scheduled workout |

---

## Training Plan Service (`/trainingplan-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/trainingplan-service/trainingplan` | List training plans |
| [ ] | GET | `/trainingplan-service/trainingplan/{planId}` | Get training plan |
| [ ] | GET | `/trainingplan-service/trainingplan/adaptive/{planId}` | Get adaptive training plan |

---

## Calendar Service (`/calendar-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/calendar-service/year/{year}` | Calendar data by year |

---

## Golf Service

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | Golf scorecard summary endpoint | Golf round summaries |
| [ ] | GET | Golf scorecard detail endpoint | Individual scorecard data |

Note: Exact golf endpoint paths are not publicly documented.

---

## Menstrual/Pregnancy Service (`/periodichealth-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/periodichealth-service/menstrualcycle/dayview/{date}` | Menstrual day view |
| [ ] | GET | `/periodichealth-service/menstrualcycle/calendar?startDate={start}&endDate={end}` | Menstrual calendar |
| [ ] | GET | `/periodichealth-service/menstrualcycle/pregnancysnapshot` | Pregnancy snapshot |

---

## Lifestyle Logging Service (`/lifestylelogging-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/lifestylelogging-service/dailyLog/{date}` | Daily log |

---

## Mobile Gateway (`/mobile-gateway/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/mobile-gateway/usersummary/trainingstatus/latest/{date}` | Latest training status |
| [ ] | GET | `/mobile-gateway/usersummary/trainingstatus/monthly/{start}/{end}` | Monthly training status |
| [ ] | GET | `/mobile-gateway/usersummary/trainingstatus/weekly/{start}/{end}` | Weekly training status |
| [ ] | GET | `/mobile-gateway/heartRate/forDate/{date}` | Heart rate for date |

---

## GraphQL Gateway

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | POST | `/graphql-gateway/graphql` | GraphQL queries |

---

## Notes

1. All endpoints require authentication via OAuth2 Bearer token
2. Dates are typically in `YYYY-MM-DD` format
3. Some endpoints require `displayName` (username) as path parameter
4. Base URL is `https://connectapi.garmin.com` (or `https://connectapi.garmin.cn` for China)
5. Some endpoints use query parameters, others use path parameters
6. The `DI-Backend` header may be required for some endpoints

## Reference Projects

| Project | Language | URL |
|---------|----------|-----|
| python-garminconnect | Python | https://github.com/cyberjunky/python-garminconnect |
| garth | Python | https://github.com/matin/garth |
| garmin-workouts | Python | https://github.com/mkuthan/garmin-workouts |
| garmin-connect | JavaScript | https://github.com/Pythe1337N/garmin-connect |
| dotnet.garmin.connect | C# | https://github.com/sealbro/dotnet.garmin.connect |
| garmy | Python | https://github.com/bes-dev/garmy |

## Priority for Implementation

### High priority (common use cases):
1. Activities - search, get, download
2. Weight - daily and range
3. Heart rate - daily
4. HRV - daily
5. User profile
6. Devices

### Medium priority:
1. Metrics (training readiness, endurance score, VO2 max)
2. Workouts (list, create, schedule)
3. Goals
4. Gear
5. User summary

### Low priority:
1. Badges/Challenges
2. Blood pressure
3. Menstrual/Pregnancy
4. Fitness age
5. Golf
6. GraphQL
