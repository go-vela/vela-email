// SPDX-License-Identifier: Apache-2.0

package main

// the default subject returns the full repository name (org/repo), the branch and the build commit.
const DefaultSubject = `
{{ .VELA_REPO_FULL_NAME }} {{ .VELA_BUILD_BRANCH }} - {{ .VELA_BUILD_COMMIT }}
`

// default html body returns the build link and build number,
// full repository name (org/repo), build author and email,
// branch, build commit, build start time, and build commit message.
const DefaultHTMLBody = `
<table>
   <tbody>
      <tr>
         <td width="600">
            <div>
               <table width="100%" cellspacing="0" cellpadding="0">
                  <tbody>
                     <tr>
                        <td>
                           <table width="100%" cellspacing="0" cellpadding="0">
                              <tbody>
                                 <tr>
                                    <td>Build Number:</td>
                                    <td><a href="{{ .VELA_BUILD_LINK }}"> 
                                    {{ .VELA_BUILD_NUMBER }} </a></td>
                                 </tr>
                                 <tr>
                                    <td>Repo:</td>
                                    <td>{{ .VELA_REPO_FULL_NAME }}</td>
                                 </tr>
                                 <tr>
                                    <td>Author:</td>
                                    <td>{{ .VELA_BUILD_AUTHOR }}
                                     ({{ .VELA_BUILD_AUTHOR_EMAIL }})</td>
                                 </tr>
                                 <tr>
                                    <td>Branch:</td>
                                    <td>{{ .VELA_BUILD_BRANCH }}</td>
                                 </tr>
                                 <tr>
                                    <td>Commit:</td>
                                    <td>{{ .VELA_BUILD_COMMIT }}</td>
                                 </tr>
                                 <tr>
                                    <td>Started at:</td>
                                    <td>{{ .BuildCreated }}</td>
                                 </tr>
                              </tbody>
                           </table>
                           <hr />
                           <table width="100%" cellspacing="0" cellpadding="0">
                              <tbody>
                                 <tr>
                                    <td>{{ .VELA_BUILD_MESSAGE }}</td>
                                 </tr>
                              </tbody>
                           </table>
                        </td>
                     </tr>
                  </tbody>
               </table>
            </div>
         </td>
      </tr>
   </tbody>
</table>
`
