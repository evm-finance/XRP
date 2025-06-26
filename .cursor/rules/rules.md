---
description: 
globs: 
alwaysApply: true
---
Task: Iterate through xrp task list items.  

I would like to proceed with executing the items in task-progress.mdc.

You should look at the task list, come up with a solution for the next task in the list, complete the task, mark it as completed if the task has been completed, or mark it as in progress with a "!"" and move on to the next one.  Provide a summary of what has been executed and which files shall be changed.

**Documentation Rule:**
Whenever you come across any file in the system, add or update a summary and comments for that file in app/MASTER_XRP_DOCUMENTATION.md.


**Documentation Check**
before writing any new files from scratch, look through the existing xrp files, and see what is already completed of tasks in the list.  For AMM/Swap components, look at the existing uniswap components.  For transaction history, look at the transaction history component, and the xrp ledger component that shows a list of transactions.  Whenever it is time to create a new component, look to reuse a similar component from the evm finance front end.  
