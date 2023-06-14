#!/usr/bin/env bash
openapi-generator-cli generate -g typescript-fetch \
	-i ./swimlogsApi.yaml \
	-o src/generated \
	--enable-post-process-file \
	--inline-schema-name-mappings getAllSessions_200_response=GetAllSessionsResponse,getTrainings_200_response=GetTrainingsResponse,getTrainingsDetails_200_response=GetDetailsResponse,getTrainingsDetailsCurrentWeek_200_response=GetDetailsCurrWeekResponse,InvalidSet_startingRule=InvalidSetStartingRule
