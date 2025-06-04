import Description from '@/components/tasks/partials/Description.vue';
import AttachmentModel from '@/models/attachment'
import type {IAttachment} from '@/modelTypes/IAttachment'

import AttachmentService from '@/services/attachment'
import {useTaskStore} from '@/stores/tasks'

import { z } from 'zod';

export const SummarySimpleSchema = z.object({
  id: z.string().optional().describe('a useless id created by AI for its own amusement'),
  summary: z.string().min(10, "Summary must be at least 10 characters long.").max(6000, "Summary must be at most 6000 characters long.").describe('Summary of the request.'),
}).describe('Summary of the document.')

export type SummarySimple = z.infer<typeof SummarySimpleSchema>
  
export const SummaryBriefSchema = z.object({
  documentId: z.string().min(2, 'The document number must be at least 2 characters long').max(64, 'The id of the object must be at most 64 characters long').or(z.literal('')).describe('The ID or Number of the document/request.'),
  requestType: z.string().optional().describe('The type of request.'),
  shortSummary: z.string().min(10, "Summary must be at least 10 characters long.").max(6000, "Summary must be at most 6000 characters long.").describe('Summary of the request.'),
  keyStakeholders: z.array(z.object({
	name: z.string().optional().describe('Name of the stakeholder.'),
	role: z.string().optional().describe('Role of the stakeholder.'),
	email: z.string().email().optional().describe('Email of the stakeholder.'),
	phone: z.string().optional().describe('Phone number of the stakeholder.'),
  })).optional().describe('Key stakeholders.'),
  importantDates: z.array(z.object({
	date: z.string().optional().describe('Date of the request.'),
	description: z.string().optional().describe('Description of the date.'),
  })).optional().describe('Important dates.'),
}).describe('Summary of the request.')

export type SummaryBrief = z.infer<typeof SummaryBriefSchema>

export async function uploadFile(taskId: number, file: File, needAi: boolean, onSuccess?: (url: string, summary: SummaryBrief) => void) {
	const attachmentService = new AttachmentService()
	const files = [file]

	return await uploadFiles(attachmentService, taskId, files, needAi, onSuccess)
}

async function summarizeFileWithAI(file: File): Promise<SummaryBrief> {

	const formData = new FormData();
	formData.append('file', file);

	const response = await fetch('/aiapi/summarizeDocument', {
		method: 'POST',
		body: formData,
	});

	if (!response.ok) {
		throw new Error(`Failed to summarize file: ${response.statusText}`);
	}

	const data = await response.json();

	// console.log('Summary response:', data);

	return SummarySimpleSchema.parse(data);
}

export async function uploadFiles(
	attachmentService: AttachmentService,
	taskId: number,
	files: File[] | FileList,
	needAi: boolean,
	onSuccess?: (attachmentUrl: string, summary: SummaryBrief) => void,
) {
	const attachmentModel = new AttachmentModel({taskId})
	const response = await attachmentService.create(attachmentModel, files)
	console.debug(`Uploaded attachments for task ${taskId}, response was`, response)

	const aiResponse = needAi ? await summarizeFileWithAI(files[0]) : null

	const task = useTaskStore()
	
	response.success?.map((attachment: IAttachment) => {
		task.addTaskAttachment({
			taskId,
			attachment,
		})
		onSuccess?.(generateAttachmentUrl(taskId, attachment.id),
		 aiResponse)
	})

	if (response.errors !== null) {
		throw Error(response.errors)
	}
}

export function generateAttachmentUrl(taskId: number, attachmentId: number) {
	return `${window.API_URL}/tasks/${taskId}/attachments/${attachmentId}`
}