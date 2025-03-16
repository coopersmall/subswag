import * as github from '@actions/github';
import * as core from '@actions/core';

const checkName = "Label Verification";
const requiredLabels = ['bug', 'enhancement'];


async function run({ context, octokit }, sha, labels) {
  const check = await octokit.rest.checks.create({
    owner: context.repo.owner,
    repo: context.repo.repo,
    name: checkName,
    head_sha: sha,
    status: 'in_progress',
  });
  
  const matchingLabels = requiredLabels.filter(required => 
    labels.includes(required.toLowerCase())
  );
  
  if (matchingLabels.length === 0) {
    await octokit.rest.checks.update({
      owner: context.repo.owner,
      repo: context.repo.repo,
      check_run_id: check.data.id,
      status: 'completed',
      conclusion: 'failure',
      output: {
        title: 'Missing Required Label',
        summary: `PR must have one of these labels: ${requiredLabels.join(', ')}`,
        text: 'Please add one of the required labels to categorize this PR appropriately.'
      }
    });
    return;
  }
   
  if (matchingLabels.length > 1) {
    await octokit.rest.checks.update({
      owner: context.repo.owner,
      repo: context.repo.repo,
      check_run_id: check.data.id,
      status: 'completed',
      conclusion: 'failure',
      output: {
        title: 'Multiple Conflicting Labels',
        summary: `PR cannot have multiple labels from: ${requiredLabels.join(', ')}`,
        text: `Current matching labels: ${matchingLabels.join(', ')}. Please remove all but one of these labels.`
      }
    });
    return;
  }
  
  await octokit.rest.checks.update({
    owner: context.repo.owner,
    repo: context.repo.repo,
    check_run_id: check.data.id,
    status: 'completed',
    conclusion: 'success',
    output: {
      title: 'Labels Verified',
      summary: `PR has the correct label: ${matchingLabels[0]}`,
      text: 'The PR is properly labeled with exactly one of the required labels.'
    }
  });
}

const sha = core.getInput('sha');
const prLabels = JSON.parse(core.getInput('labels')).map(label =>
  label.name.toLowerCase()
);
const context = github.context;
const octokit = github.getOctokit(core.getInput('token'));

await run({context, octokit}, sha, prLabels);
