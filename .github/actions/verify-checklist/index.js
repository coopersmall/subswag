import * as github from '@actions/github';
import * as core from '@actions/core';

const checkName = "Checklist Verification";

const issueNumber = core.getInput('issue_number');
 
const response = await github.rest.issues.listComments({
    owner: github.context.repo.owner,
    repo: github.context.repo.repo,
    issue_number: issueNumber
});
 
const checklistComment = response.data.find(comment => 
    comment.body.includes('## Required Acknowledgements')
);
 
const check = await github.rest.checks.create({
    owner: github.context.repo.owner,
    repo: github.context.repo.repo,
    name: checkName,
    head_sha: github.context.sha,
    status: 'in_progress',
});
 
if (!checklistComment) {
    await github.rest.checks.update({
        owner: github.context.repo.owner,
        repo: github.context.repo.repo,
        check_run_id: check.data.id,
        status: 'completed',
        conclusion: 'failure',
        output: {
            title: 'Required Checklist Missing',
            summary: 'Could not find the required acknowledgements checklist in the PR comments.',
            text: 'Please ensure that the PR includes the Required Acknowledgements checklist.'
        }
    });
} else {
  const uncheckedBoxes = (checklistComment.body.match(/\[ \]/g) || []).length;
  if (uncheckedBoxes > 0) {
      await github.rest.checks.update({
          owner: github.context.repo.owner,
          repo: github.context.repo.repo,
          check_run_id: check.data.id,
          status: 'completed',
          conclusion: 'failure',
          output: {
              title: 'Missing Checklist Items',
              summary: `Found ${uncheckedBoxes} unchecked item${uncheckedBoxes > 1 ? 's' : ''} in the PR checklist`,
              text: 'Please complete all items in the Required Acknowledgements checklist before proceeding.'
          }
      });
  } else {
      await github.rest.checks.update({
          owner: github.context.repo.owner,
          repo: github.context.repo.repo,
          check_run_id: check.data.id,
          status: 'completed',
          conclusion: 'success',
          output: {
              title: 'Checklist Complete',
              summary: 'All required checklist items have been completed.',
              text: 'The PR checklist has been properly filled out.'
          }
      });
  }
}
