#!/usr/bin/env bash

openapi-generator-cli generate -g typescript-fetch -i ./swimlogsAPI.gen.yaml \
	--inline-schema-name-mappings summariesCurrentWeek_200_response=SummariesCurrentWeekResponse \
	--inline-schema-name-mappings summariesPage_200_response=SummariesPageResponse \
	--inline-schema-name-mappings moveSet_request=MoveSetRequest \
	--inline-schema-name-mappings editSet_request=EditSetRequest \
	--inline-schema-name-mappings editSet_200_response=EditSetResponse \
	--inline-schema-name-mappings trainingById_404_response=TrainingNotFoundResponse \
	--inline-schema-name-mappings createTraining_400_response=InvalidTrainingResponse
