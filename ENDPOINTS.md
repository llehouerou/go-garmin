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
| [x] | GET | `/wellness-service/wellness/daily/spo2/{date}` | Daily SpO2 data |
| [x] | GET | `/wellness-service/wellness/daily/respiration/{date}` | Daily respiration data |
| [x] | GET | `/wellness-service/wellness/daily/im/{date}` | Daily intensity minutes |
| [ ] | GET | `/wellness-service/wellness/dailyEvents/{date}` | Daily events |
| [ ] | GET | `/wellness-service/wellness/dailySleepData/{displayName}?date={date}` | Daily sleep (alternative) |
| [ ] | GET | `/wellness-service/wellness/dailySummaryChart/{displayName}?date={date}` | Daily summary chart (steps) |
| [ ] | GET | `/wellness-service/wellness/floorsChartData/daily/{date}` | Floor climbing data |
| [ ] | POST | `/wellness-service/wellness/epoch/request/{date}` | Request epoch data reload |
| [ ] | GET | `/wellness-service/wellness/bodyBattery/reports/daily?startDate={start}&endDate={end}` | Body battery reports |
| [ ] | GET | `/wellness-service/stats/daily/sleep/score/{start}/{end}` | Sleep score stats |

Note: Some endpoints like `dailyHeartRate` can also use `/{displayName}?date={date}` format.

---

## Activity List Service (`/activitylist-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/activitylist-service/activities/search/activities?start={start}&limit={limit}` | Search activities |
| [ ] | GET | `/activitylist-service/activities/` | List activities |
| [ ] | GET | `/activitylist-service/activities/count` | Activity count |
| [ ] | GET | `/activitylist-service/activities/{gearUUID}/gear?start={start}&limit={limit}` | Activities for gear |

---

## Activity Service (`/activity-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/activity-service/activity/{activityId}` | Get single activity |
| [x] | GET | `/activity-service/activity/{activityId}/details?maxChartSize={n}&maxPolylineSize={n}` | Activity details (time-series) |
| [x] | GET | `/activity-service/activity/{activityId}/splits` | Activity splits |
| [ ] | GET | `/activity-service/activity/{activityId}/typedsplits` | Activity typed splits |
| [ ] | GET | `/activity-service/activity/{activityId}/split_summaries` | Activity split summaries |
| [x] | GET | `/activity-service/activity/{activityId}/weather` | Activity weather |
| [x] | GET | `/activity-service/activity/{activityId}/hrTimeInZones` | HR time in zones |
| [x] | GET | `/activity-service/activity/{activityId}/powerTimeInZones` | Power time in zones |
| [x] | GET | `/activity-service/activity/{activityId}/exerciseSets` | Exercise sets |
| [ ] | GET | `/activity-service/activity/{activityId}/gear` | Activity gear |
| [ ] | GET | `/activity-service/activity/activityTypes` | Activity types |
| [ ] | POST | `/activity-service/activity` | Create manual activity |
| [ ] | PUT | `/activity-service/activity/{activityId}` | Update activity (name, type, etc.) |
| [ ] | DELETE | `/activity-service/activity/{activityId}` | Delete activity |

Note: `typedsplits` is lowercase 'd' in the actual API.

---

## Download Service (`/download-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/download-service/files/activity/{activityId}` | Download activity (original FIT) |
| [x] | GET | `/download-service/export/tcx/activity/{activityId}` | Export activity as TCX |
| [x] | GET | `/download-service/export/gpx/activity/{activityId}` | Export activity as GPX |
| [x] | GET | `/download-service/export/kml/activity/{activityId}` | Export activity as KML |
| [x] | GET | `/download-service/export/csv/activity/{activityId}` | Export activity as CSV |

---

## Upload Service (`/upload-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | POST | `/upload-service/upload` | Upload activity file (FIT, TCX, GPX) |

---

## Weight Service (`/weight-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/weight-service/weight/dayview/{date}` | Daily weight data |
| [x] | GET | `/weight-service/weight/range/{start}/{end}?includeAll=true` | Weight range |
| [ ] | GET | `/weight-service/weight/dateRange?startDate={start}&endDate={end}` | Weight date range |
| [ ] | GET | `/weight-service/weight/daterangesnapshot` | Body composition snapshot |
| [ ] | POST | `/weight-service/user-weight` | Add weigh-in |
| [ ] | DELETE | `/weight-service/weight/{date}/byversion/{weightPK}` | Delete weigh-in |

---

## User Summary Service (`/usersummary-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/usersummary-service/usersummary/daily/{displayName}?calendarDate={date}` | Daily user summary |
| [ ] | GET | `/usersummary-service/usersummary/hydration/daily/{date}` | Daily hydration |
| [ ] | POST | `/usersummary-service/usersummary/hydration/log` | Log/update hydration |
| [ ] | GET | `/usersummary-service/stats/steps/daily/{start}/{end}` | Daily steps stats (max 28 days) |
| [ ] | GET | `/usersummary-service/stats/steps/weekly/{end}/{weeks}` | Weekly steps stats |
| [ ] | GET | `/usersummary-service/stats/stress/daily/{start}/{end}` | Daily stress stats |
| [ ] | GET | `/usersummary-service/stats/stress/weekly/{end}/{weeks}` | Weekly stress stats |
| [ ] | GET | `/usersummary-service/stats/hydration/daily/{start}/{end}` | Hydration stats |
| [ ] | GET | `/usersummary-service/stats/im/daily/{start}/{end}` | Daily intensity minutes |
| [ ] | GET | `/usersummary-service/stats/im/weekly/{start}/{end}` | Weekly intensity minutes |

Note: `daily/{displayName}` can also be accessed as `daily/?calendarDate={date}` (garth variant).

---

## User Stats Service (`/userstats-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/userstats-service/wellness/daily/{displayName}?fromDate={date}&untilDate={date}&metricId=60` | Daily wellness stats (RHR) |

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
| [x] | GET | `/metrics-service/metrics/trainingreadiness/{date}` | Training readiness |
| [x] | GET | `/metrics-service/metrics/endurancescore?calendarDate={date}` | Endurance score |
| [ ] | GET | `/metrics-service/metrics/endurancescore/stats?startDate={start}&endDate={end}&aggregation={agg}` | Endurance score stats |
| [x] | GET | `/metrics-service/metrics/hillscore?calendarDate={date}` | Hill score |
| [ ] | GET | `/metrics-service/metrics/hillscore/stats?startDate={start}&endDate={end}&aggregation={agg}` | Hill score stats |
| [ ] | GET | `/metrics-service/metrics/racepredictions/latest/{displayName}` | Latest race predictions (requires display name) |
| [ ] | GET | `/metrics-service/metrics/racepredictions/daily/{displayName}?_={timestamp}` | Daily race predictions |
| [ ] | GET | `/metrics-service/metrics/racepredictions/monthly/{displayName}?_={timestamp}` | Monthly race predictions |
| [x] | GET | `/metrics-service/metrics/maxmet/daily/{start}/{end}` | Daily VO2 max/MET |
| [x] | GET | `/metrics-service/metrics/maxmet/latest/{date}` | Latest VO2 max/MET |
| [x] | GET | `/metrics-service/metrics/trainingstatus/aggregated/{date}` | Training status aggregated |
| [x] | GET | `/metrics-service/metrics/trainingstatus/daily/{date}` | Daily training status |
| [x] | GET | `/metrics-service/metrics/trainingloadbalance/latest/{date}` | Training load balance |
| [x] | GET | `/metrics-service/metrics/heataltitudeacclimation/latest/{date}` | Heat/altitude acclimation |

---

## Device Service (`/device-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/device-service/deviceregistration/devices` | List devices |
| [x] | GET | `/device-service/deviceservice/device-info/settings/{deviceId}` | Device settings |
| [x] | GET | `/device-service/devicemessage/messages` | Device messages |
| [ ] | GET | `/device-service/deviceservice/mylastused` | Last used device |

---

## Web Gateway (`/web-gateway/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/web-gateway/device-info/primary-training-device` | Primary training device |
| [ ] | GET | `/web-gateway/solar/{deviceId}/{startDate}/{endDate}` | Solar panel data |

---

## User Profile Service (`/userprofile-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/userprofile-service/socialProfile` | Social profile (displayName, fullName) |
| [x] | GET | `/userprofile-service/userprofile/user-settings` | User settings (measurementSystem) |
| [x] | GET | `/userprofile-service/userprofile/settings` | Profile settings |
| [ ] | GET | `/userprofile-service/userprofile/profile` | User profile details |

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
| [ ] | GET | `/gear-service/gear/user/{userProfilePk}/activityTypes` | Gear activity types |
| [ ] | POST | `/gear-service/gear/link/{gearUUID}/activity/{activityId}` | Link gear to activity |
| [ ] | POST | `/gear-service/gear/unlink/{gearUUID}/activity/{activityId}` | Unlink gear from activity |

Note: Gear activities are fetched via `/activitylist-service/activities/{gearUUID}/gear`

---

## Badge Service (`/badge-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/badge-service/badge/earned` | Earned badges |
| [ ] | GET | `/badge-service/badge/available?showExclusiveBadge=true` | Available badges |

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
| [ ] | DELETE | `/bloodpressure-service/bloodpressure/{date}/{version}` | Delete blood pressure |

---

## Personal Record Service (`/personalrecord-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/personalrecord-service/personalrecord/prs/{displayName}` | Personal records (requires display name) |

---

## Biometric Service (`/biometric-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/biometric-service/biometric/latestLactateThreshold` | Latest lactate threshold |
| [x] | GET | `/biometric-service/biometric/latestFunctionalThresholdPower/CYCLING` | Latest cycling FTP |
| [x] | GET | `/biometric-service/biometric/powerToWeight/latest/{date}?sport=Running` | Power-to-weight ratio |
| [x] | GET | `/biometric-service/stats/lactateThresholdSpeed/range/{start}/{end}?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST` | LT speed range |
| [x] | GET | `/biometric-service/stats/lactateThresholdHeartRate/range/{start}/{end}?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST` | LT heart rate range |
| [x] | GET | `/biometric-service/stats/functionalThresholdPower/range/{start}/{end}?sport=RUNNING&aggregation=daily&aggregationStrategy=LATEST` | FTP range |
| [x] | GET | `/biometric-service/heartRateZones/` | Heart rate zones for all sports |

---

## Fitness Age Service (`/fitnessage-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/fitnessage-service/fitnessage/{date}` | Fitness age |

---

## Fitness Stats Service (`/fitnessstats-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/fitnessstats-service/activity/all?startDate={start}&endDate={end}&...` | All activity stats |

---

## Workout Service (`/workout-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [x] | GET | `/workout-service/workouts?start={start}&limit={limit}` | List workouts |
| [x] | GET | `/workout-service/workout/{workoutId}` | Get workout |
| [x] | GET | `/workout-service/workout/FIT/{workoutId}` | Download workout as FIT |
| [x] | POST | `/workout-service/workout` | Create workout |
| [x] | PUT | `/workout-service/workout/{workoutId}` | Update workout |
| [x] | DELETE | `/workout-service/workout/{workoutId}` | Delete workout |
| [x] | POST | `/workout-service/schedule/{workoutId}` | Schedule workout (body: {"date": "YYYY-MM-DD"}) |
| [x] | GET | `/workout-service/schedule/{scheduleId}` | Get scheduled workout |

Note: Can also be accessed via `/proxy/workout-service/` prefix (garmin-workouts).

---

## Training Plan Service (`/trainingplan-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/trainingplan-service/trainingplan/plans` | List training plans |
| [ ] | GET | `/trainingplan-service/trainingplan/phased/{planId}` | Get phased training plan |
| [ ] | GET | `/trainingplan-service/trainingplan/fbt-adaptive/{planId}` | Get FBT adaptive plan |

---

## Calendar Service (`/calendar-service/`)

| Status | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| [ ] | GET | `/calendar-service/year/{year}` | Calendar data by year (unverified) |

Note: This endpoint could not be verified in any reference implementation.

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
| [ ] | GET | `/periodichealth-service/menstrualcycle/calendar/{start}/{end}` | Menstrual calendar |
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
| [ ] | GET | `/mobile-gateway/heartRate/forDate/{date}` | Heart rate for date |
| [ ] | GET | `/mobile-gateway/usersummary/trainingstatus/latest/{date}` | Latest training status (unverified) |
| [ ] | GET | `/mobile-gateway/usersummary/trainingstatus/monthly/{start}/{end}` | Monthly training status (unverified) |
| [ ] | GET | `/mobile-gateway/usersummary/trainingstatus/weekly/{start}/{end}` | Weekly training status (unverified) |

Note: heartRate/forDate verified in python-garminconnect. Training status endpoints may be mobile app specific.

---

## GraphQL Gateway (`/graphql-gateway/`)

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
