#!/usr/bin/env bash

openapi-generator-cli generate \
	-g typescript-fetch \
	-i ../api/swimlogsAPI.gen.yaml \
	-o ./src/api/generated/ \
	--inline-schema-name-mappings replaceSetComponents_request=ReplaceSetComponentsRequest \
	--inline-schema-name-mappings editSet_request=EditSetRequest \
	--inline-schema-name-mappings moveSet_request=MoveSetRequest \
	--inline-schema-name-mappings summariesPage_200_response=SummariesPageResponse
