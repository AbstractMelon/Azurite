import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const API_BASE_URL = 'http://localhost:8080';

export const GET: RequestHandler = async ({ params, fetch }) => {
	const { gameSlug, modSlug } = params;

	try {
		const response = await fetch(`${API_BASE_URL}/download/${gameSlug}/${modSlug}`);

		if (!response.ok) {
			throw error(response.status, 'Download failed');
		}

		// Get response as blob
		const blob = await response.blob();

		// Forward important headers
		const headers = new Headers();
		const contentType = response.headers.get('content-type');
		const contentDisposition = response.headers.get('content-disposition');

		if (contentType) headers.set('Content-Type', contentType);
		if (contentDisposition) headers.set('Content-Disposition', contentDisposition);

		return new Response(blob, {
			status: 200,
			headers
		});
	} catch (err) {
		console.error('Download error:', err);
		throw error(500, 'Failed to download file');
	}
};
