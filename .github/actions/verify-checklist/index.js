import * as github from '@actions/github';
import * as core from '@actions/core';

const checkName = "Checklist Verification";

async function run({context, octokit}, issueNumber) {
  const response = await octokit.rest.issues.listComments({
      owner: context.repo.owner,
      repo: context.repo.repo,
      issue_number: issueNumber
  });
   
  const checklistComment = response.data.find(comment => 
      comment.body.includes('## Required Acknowledgements')
  );
   
  const check = await octokit.rest.checks.create({
      owner: context.repo.owner,
      repo: context.repo.repo,
      name: checkName,
      head_sha: context.sha,
      status: 'in_progress',
  });
   
  if (!checklistComment) {
      await octokit.rest.checks.update({
          owner: context.repo.owner,
          repo: context.repo.repo,
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
        await octokit.rest.checks.update({
            owner: context.repo.owner,
            repo: context.repo.repo,
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
        await octokit.rest.checks.update({
            owner: context.repo.owner,
            repo: context.repo.repo,
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
}


const issueNumber = core.getInput('issue_number');
const context = github.context;
const octokit = github.getOctokit(core.getInput('token'));

run({context, octokit}, issueNumber);
