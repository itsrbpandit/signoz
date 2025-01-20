/* eslint-disable sonarjs/no-duplicate-string */
import { Pipeline } from 'types/api/pipeline/def';
import { DataTypes } from 'types/api/queryBuilder/queryAutocompleteResponse';

export const pipelineData: Pipeline = {
	id: 'test-id-1',
	version: 24,
	elementType: 'log_pipelines',
	active: false,
	is_valid: false,
	disabled: false,
	deployStatus: 'DEPLOYED',
	deployResult: 'Deployment was successful',
	lastHash: 'log_pipelines:24',
	lastConf: 'oiwernveroi',
	createdBy: 'test-created-by',
	pipelines: [
		{
			id: 'test-id-2',
			orderId: 1,
			name: 'hotrod logs parser',
			alias: 'hotrodlogsparser',
			description: 'Trying to test Logs Pipeline feature',
			enabled: true,
			filter: {
				op: 'AND',
				items: [
					{
						key: {
							key: 'container_name',
							dataType: DataTypes.String,
							type: 'tag',
							isColumn: false,
							isJSON: false,
						},
						id: 'sampleid',
						value: 'hotrod',
						op: '=',
					},
				],
			},
			config: [
				{
					type: 'regex_parser',
					id: 'parsetext(regex)',
					output: 'parseattribsjson',
					on_error: 'send',
					orderId: 1,
					enabled: true,
					name: 'parse text (regex)',
					parse_to: 'attributes',
					regex:
						'.+\\t+(?P<log_level>.+)\\t+(?P<location>.+)\\t+(?P<message>.+)\\t+(?P<attribs_json>.+)',
					parse_from: 'body',
				},
				{
					type: 'json_parser',
					id: 'parseattribsjson',
					output: 'removetempattribs_json',
					orderId: 2,
					enabled: true,
					name: 'parse attribs json',
					parse_to: 'attributes',
					parse_from: 'attributes.attribs_json',
				},
				{
					type: 'remove',
					id: 'removetempattribs_json',
					output: 'c2062723-895e-4614-ba38-29c5d5ee5927',
					orderId: 3,
					enabled: true,
					name: 'remove temp attribs_json',
					field: 'attributes.attribs_json',
				},
				{
					type: 'add',
					id: 'c2062723-895e-4614-ba38-29c5d5ee5927',
					orderId: 4,
					enabled: true,
					name: 'test add ',
					field: 'resource["container.name"]',
					value: 'hotrod',
				},
			],
			createdBy: 'test@email',
			createdAt: '2024-01-02T13:56:02.858300964Z',
		},
		{
			id: 'tes-id-1',
			orderId: 2,
			name: 'Logs Parser - test - Customer Service',
			alias: 'LogsParser-test-CustomerService',
			description: 'Trying to test Logs Pipeline feature',
			enabled: true,
			filter: {
				op: 'AND',
				items: [
					{
						key: {
							key: 'service',
							dataType: DataTypes.String,
							type: 'tag',
							isColumn: false,
							isJSON: false,
						},
						id: 'sample-test-1',
						value: 'customer',
						op: '=',
					},
				],
			},
			config: [
				{
					type: 'grok_parser',
					id: 'Testtest',
					on_error: 'send',
					orderId: 1,
					enabled: true,
					name: 'Test test',
					parse_to: 'attributes',
					pattern:
						'^%{DATE:date}Z INFO customer/database.go:73 Loading customer {"service": "customer", "component": "mysql", "trace_id": "test-id", "span_id": "1427a3fcad8b1514", "customer_id": "567"}',
					parse_from: 'body',
				},
			],
			createdBy: 'test@email',
			createdAt: '2024-01-02T13:56:02.863764227Z',
		},
	],
	history: [
		{
			id: 'test-id-4',
			version: 24,
			elementType: 'log_pipelines',
			active: false,
			isValid: false,
			disabled: false,
			deployStatus: 'DEPLOYED',
			deployResult: 'Deployment was successful',
			lastHash: 'log_pipelines:24',
			lastConf: 'eovineroiv',
			createdBy: 'test-created-by',
			createdByName: 'test-user',
			createdAt: '2024-01-02T13:56:02Z',
		},
		{
			id: 'test-4',
			version: 23,
			elementType: 'log_pipelines',
			active: false,
			isValid: false,
			disabled: false,
			deployStatus: 'DEPLOYED',
			deployResult: 'Deployment was successful',
			lastHash: 'log_pipelines:23',
			lastConf: 'eivrounreovi',
			createdBy: 'test-created-by',
			createdByName: 'test-user',
			createdAt: '2023-12-29T12:59:20Z',
		},
	],
};

export const pipelineDataHistory: Pipeline['history'] = [
	{
		id: 'test-id-4',
		version: 24,
		elementType: 'log_pipelines',
		active: false,
		isValid: false,
		disabled: false,
		deployStatus: 'DEPLOYED',
		deployResult: 'Deployment was successful',
		lastHash: 'log_pipelines:24',
		lastConf: 'eovineroiv',
		createdBy: 'test-created-by',
		createdByName: 'test-user',
		createdAt: '2024-01-02T13:56:02Z',
	},
	{
		id: 'test-4',
		version: 23,
		elementType: 'log_pipelines',
		active: false,
		isValid: false,
		disabled: false,
		deployStatus: 'IN_PROGRESS',
		deployResult: 'Deployment is in progress',
		lastHash: 'log_pipelines:23',
		lastConf: 'eivrounreovi',
		createdBy: 'test-created-by',
		createdByName: 'test-user',
		createdAt: '2023-12-29T12:59:20Z',
	},
	{
		id: 'test-4-1',
		version: 25,
		elementType: 'log_pipelines',
		active: false,
		isValid: false,
		disabled: false,
		deployStatus: 'DIRTY',
		deployResult: 'Deployment is dirty',
		lastHash: 'log_pipelines:23',
		lastConf: 'eivrounreovi',
		createdBy: 'test-created-by',
		createdByName: 'test-user',
		createdAt: '2023-12-29T12:59:20Z',
	},
	{
		id: 'test-4-2',
		version: 26,
		elementType: 'log_pipelines',
		active: false,
		isValid: false,
		disabled: false,
		deployStatus: 'FAILED',
		deployResult: 'Deployment failed',
		lastHash: 'log_pipelines:23',
		lastConf: 'eivrounreovi',
		createdBy: 'test-created-by',
		createdByName: 'test-user',
		createdAt: '2023-12-29T12:59:20Z',
	},
	{
		id: 'test-4-3',
		version: 27,
		elementType: 'log_pipelines',
		active: false,
		isValid: false,
		disabled: false,
		deployStatus: 'UNKNOWN',
		deployResult: '',
		lastHash: 'log_pipelines:23',
		lastConf: 'eivrounreovi',
		createdBy: 'test-created-by',
		createdByName: 'test-user',
		createdAt: '2023-12-29T12:59:20Z',
	},
];
