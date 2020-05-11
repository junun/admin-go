import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Row,
  Col, Button, Popconfirm, Icon, message, Select} from "antd";
import { Link } from 'react-router-dom';
import {connect} from "dva";
import {timeTrans, timeDatetimeTrans, hasPermission} from "@/utils/globalTools"

import WsDeployMessage from "@/components/Report/WsDeployMessage";
import WsUndoMessage from "@/components/Report/WsUndoMessage";

const FormItem = Form.Item;
const Option = Select.Option;

const checkStatus = {
  '1': '新建上线单',
  '2': '审核成功',
  '3': '审核失败',
  '4': '上线失败',
  '5': '上线成功',
  '6': '回滚成功',
  '7': '回滚失败',
};

@connect(({ loading, deploy, config }) => {
    return {
      deployList: deploy.deployList,
      deployListLoading: loading.effects['deploy/getDeploy'],
      deployLen: deploy.deployLen,
      apiToken: deploy.apiToken,
      projectsList: config.projectsList,
      appTemplateList: config.appTemplateList,
      configEnvList: config.configEnvList,
      gitBranchList: deploy.gitBranchList,
      gitCommitList: deploy.gitCommitList,
    }
 })

class DeployListPage extends React.Component {
  state = {
    addModal: false,
    visible: false,
    appVisible: false,
    editCacheData: {},
    appId: 0,
    gitShow: 0,
    disableValues: false,
    showTemplate: false,
    undoTemplate: false,
    id: 0,
    branchLoading: false,
    commitLoading: false,
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'deploy/getDeploy',
      payload: {
        page: 1,
        pageSize: 10,
      }
    });

    dispatch({
      type: 'config/getConfigEnv',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });

    dispatch({
      type: 'config/getAppTemplate',
      payload: {
        aid: 0,
      }
    });
  }
  
  showAddModal = () => {
    this.setState({ 
      appVisible: true,
      addModal: true,
    });
  };

  handleCancel = () => {
    this.setState({
      appVisible: false,
      addModal: false,
    });
  };

  handleEx1Cancel = () => {
    this.setState({
      visible: false,
      gitShow: 0,
      addModal: false,
      disableValues: false,
      editCacheData: {},
    });
    // 清理 branchList
    const { dispatch } = this.props;
    dispatch({
      type: 'deploy/cleanBranchList',
    });
  };

  switchTemplate = () => {
    this.setState({ 
      showTemplate: !this.state.showTemplate,
    });
  };

  switchUndoTemplate = () => {
    this.setState({ 
      undoTemplate: !this.state.undoTemplate,
    });
    const { dispatch } = this.props;
    dispatch({ 
      type: 'deploy/getDeploy',
      payload: {
        page: 1,
        pageSize: 10,
      }
    });
  };

  handleEnvChange = (values) => {
    if (values) {
      const { dispatch } = this.props;
      dispatch({
        type: 'config/getProject',
        payload: {
          page: 1,
          pageSize: 999,
          active: 1,
          envId: values,
        }
      });
    }
  };

  handleAppChange = (values) => {
    if (values) {
      const { dispatch } = this.props;
      dispatch({
        type: 'config/getAppTemplate',
        payload: {
          aid: values,
        }
      });
    }
  };

  handleExt1Ok = () => {
    this.setState({ 
      visible: true,
      appVisible: false,
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        if (!this.state.addModal) {
          if (
            obj.Name       === values.Name && 
            obj.RepoBranch === values.RepoBranch &&
            obj.RepoCommit === values.RepoCommit 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.ID;
            values.RepoCommit = values.RepoCommit.split(" ")[0]
            dispatch({
              type: 'deploy/deployEdit',
              payload: values,
            });
          }
        } else {
          values.RepoCommit = values.RepoCommit.split(" ")[0]
          dispatch({
            type: 'deploy/deployAdd',
            payload: values,
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({
          visible: false,
          gitShow: 0,
          addModal: false,
        });
      }
    });
  };

  // 查branch内容
  handleChange = (values) => {
    if (values) {
      // 根据模板Tdid查branch
      this.setState({
        appId: values,
        gitShow: 1,
        branchLoading: true,
      });
      const { dispatch } = this.props;
      dispatch({
        type: 'deploy/getGitBranch',
        payload: values,
      }).then(()=>{
        this.setState({
          branchLoading: false,
        });
      });
    }
  };

  // 查commit 信息
  handleBrachChange = (values) => {
    if (values) {
      const editCache = this.state.editCacheData;
      // 判断发布关联方式
      const items = {};
      items.aid  = this.state.appId;
      items.name = values
      // 根据项目id查branch 和 commit号
      this.setState({
        commitLoading: true,
      });
      const { dispatch } = this.props;
      dispatch({
        type: 'deploy/getGitCommit',
        payload: items,
      }).then(()=>{
        this.setState({
          commitLoading: false,
        });
      });
      editCache.repo_branch = '';
      this.setState({
        editCacheData: editCache,
      });
    }
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'deploy/deployDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  // Popconfirm 取消事件
  cancel = () => {
  };

  // 审核拒绝
  powerCancel = (values) => {
    const { dispatch } = this.props;
    values.status = 3
    dispatch({
      type: 'deploy/deployReview',
      payload: values,
    });
  };
  
  // 审核通过
  powerOk = (values) => {
    const { dispatch } = this.props;
    values.status = 2
    dispatch({
      type: 'deploy/deployReview',
      payload: values,
    });
  };

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑-' + values.Name;
    this.setState({ 
      visible: true ,
      appId: values.Aid,
      editCacheData: values,
      disableValues: true,
      gitShow: 1,
      branchLoading: true,
      commitLoading: true,
    });
    const { dispatch } = this.props;
    dispatch({
      type: 'deploy/getGitBranch',
      payload: values.Aid,
    }).then(()=>{
      this.setState({ 
        branchLoading: false,
      });
    });

    const items = {};
    items.aid  = values.Aid;
    items.name = values.RepoBranch;
    dispatch({
      type: 'deploy/getGitCommit',
      payload: items,
    }).then(()=>{
      this.setState({ 
        commitLoading: false,
      });
    });
  };

  undeployYes = (values) => {
    this.setState({ 
      undoTemplate: true ,
      id: values,
    });
  };

  // 发布
  deployYes = (values) => {
    this.setState({ 
      showTemplate: true ,
      id: values,
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'deploy/getDeploy',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  columns = [
    {
      title: 'ID',
      dataIndex: 'ID',
    },
    {
      title: '申请标题',
      dataIndex: 'Name',
    },
    {
      title: '上线模板',
      dataIndex: 'Tid',
      'render': Tid => this.props.appTemplateList.map(x => {
        if (Tid==x.Dtid) {
          return x.TemplateName
        }
      })
    },
    {
      title: '分支',
      dataIndex: 'RepoBranch',
    },
    {
      title: '版本',
      dataIndex: 'RepoCommit',
    },
    {
      title: '状态',
      dataIndex: 'Status',
      'render': Status =>  checkStatus[Status]
    },
    ,
    {
      title: '状态变更时间',
      dataIndex: 'UpdateTime',
      'render': UpdateTime => timeDatetimeTrans(UpdateTime),
    },
    {
      title: '操作',
      key: 'action',
      render: (text, record) => (
        <span>
          {hasPermission('deploy-app-deploy') && record.Status == 2  &&
            <Popconfirm title="你确定要发布吗?" okText="发布" 
              onConfirm={()=>{this.deployYes(record.ID)}} onCancel={()=>{this.cancel()}}>
              <a title="上线" ><Icon type="redo" />上线</a>
            </Popconfirm>
          } 
          {hasPermission('deploy-app-undo') && record.Status >= 4 && record.Status < 6 &&
            <Popconfirm title="你确定要回退到上一个版本吗?" okText="回滚" 
              onConfirm={()=>{this.undeployYes(record.ID)}} onCancel={()=>{this.cancel()}}>
              <a title="回滚" ><Icon type="undo" />回滚</a>
            </Popconfirm>
          } 
          {hasPermission('deploy-app-review') && record.Status == 1  &&
            <Popconfirm title="审核是否通过?"  cancelText="驳回" okText="通过" 
              onConfirm={()=>{this.powerOk(record)}}
              onCancel={()=>{this.powerCancel(record)}}>
              <a title="删除" ><Icon type="team" />审核</a>
            </Popconfirm>
          }
          { record.Status == 1  &&  <Divider type="vertical" /> }
          {hasPermission('deploy-app-edit') && record.Status == 1  && 
            <a onClick={()=>{this.handleEdit(record)}}>
              <Icon type="edit"/>编辑
            </a>
          }
          { record.Status == 1 && <Divider type="vertical" /> }
          {hasPermission('deploy-app-del') && record.Status == 1  &&
            <Popconfirm title="你确定要删除吗?"
              onConfirm={()=>{this.deleteRecord(record.ID)}} onCancel={()=>{this.cancel()}}>
              <a title="删除" ><Icon type="delete" />删除</a>
            </Popconfirm>
          }
        </span>
      ),
    },
  ];

  render() {
    const {visible, appVisible, showTemplate, undoTemplate,  editCacheData, gitShow,
      disableValues, commitLoading, branchLoading} = this.state;

    const {deployList, deployListLoading, deployLen, projectsList,
      configEnvList, gitBranchList, gitCommitList, output, 
      appTemplateList, form: { getFieldDecorator } } = this.props;

    const addvar = <Button type="" onClick={this.showAddModal} >提单</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('deploy-app-newline') && <Col span={10}>{addvar}</Col>}
      </Row>;
    return (
      <div>

        <Modal
          title= { editCacheData.title || "选择项目信息" }
          visible= {appVisible}
          width={800}
          destroyOnClose= "true"
          okText="下一步"
          onOk={this.handleExt1Ok}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem  required label="选择环境">
              <Select
                placeholder="Please select"
                onChange={this.handleEnvChange}
                style={{ width: '100%' }} 
              >
              {configEnvList.map(x => <Option key={x.id} value={x.id}>{x.Name}</Option>)}
              </Select>
            </FormItem>
            <FormItem label="选择项目">
              
                <Select
                  placeholder="Please select"
                  onChange={this.handleAppChange}
                  style={{ width: '100%' }} 
                >
                {projectsList.map(x => <Option key={x.id} value={x.id}>{x.Name}</Option>)}
                </Select>
              
            </FormItem>
         </Form> 
        </Modal>

        <Modal
          title= { editCacheData.title || "新建上线单" }
          visible= {visible}
          width={800}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleEx1Cancel}
        >
          <Form>
            <FormItem label="上线单标题">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="择项目发布模板">
              <Col span={16}>
                {getFieldDecorator('Tid', {
                  initialValue: editCacheData.Tid || 'Please select' ,
                  rules: [{ required: true }],
                })( 
                  <Select
                    placeholder="Please select"
                    onChange={this.handleChange}
                    style={{ width: '100%' }} 
                    disabled={disableValues}
                  >
                  {appTemplateList.map(x => <Option key={x.Dtid} value={x.Dtid}>{x.TemplateName}</Option>)}
                  </Select>
                )}
              </Col>
              <Col span={6} offset={2}>
                <Link to="/config/app">新建发布模板</Link>
              </Col>
            </FormItem>
            { gitShow == 1 && 
              <FormItem label="选择分支" help="根据网络情况，刷新可能会很慢，请耐心等待。">
                {getFieldDecorator('RepoBranch', {
                  initialValue: editCacheData.RepoBranch || '',
                  rules: [{ required: true }],
                })(
                  <Select
                    placeholder="Please select"
                    onChange={this.handleBrachChange}
                    style={{ width: '100%' }}
                    loading={branchLoading}
                  >
                  {gitBranchList.map(x => <Option key={x} value={x}>{x}</Option>)}
                  </Select>
                )}
              </FormItem>
            }
            { gitShow == 1 && 
              <FormItem label="选择版本">
                {getFieldDecorator('RepoCommit', {
                  initialValue: editCacheData.RepoCommit || '',
                  rules: [{ required: true }],
                })(
                  <Select
                    placeholder="Please select"
                    // onChange={this.handleBrachChange}
                    style={{ width: '100%' }}
                    loading={commitLoading}
                  >
                  {gitCommitList.map(x => <Option key={x} value={x}>{x}</Option>)}
                  </Select>
                )}
              </FormItem>
            }
            { gitShow == 1 && 
              <FormItem wrapperCol={{span: 14, offset: 6}}>
                <span role="img" aria-label="notice">⚠️请先选择分支，再选版本。</span>
              </FormItem>
            }
          </Form> 
        </Modal>

        {showTemplate && 
          <WsDeployMessage id={this.state.id} onCancel={this.switchTemplate} onOk={body => this.setState({body})}/>
        }
        {undoTemplate && 
          <WsUndoMessage id={this.state.id} onCancel={this.switchUndoTemplate} onOk={body => this.setState({body})}/>
        }

        <Card title="" extra={extra}>
          <Table 
          pagination={{
            showQuickJumper: true,
            total: deployLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={deployList} loading={deployListLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(DeployListPage);
